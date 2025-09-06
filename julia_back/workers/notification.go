package workers

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SherClockHolmes/webpush-go"
)

type Config struct {
	BatchInterval time.Duration
	BatchSize     int
	MaxRetries    int
	VAPIDPublic   string
	VAPIDPrivate  string
	SubscriberURL string
	TTL           int
}

func DefaultConfig() *Config {
	return &Config{
		BatchInterval: 15 * time.Second,
		BatchSize:     100,
		MaxRetries:    5,
		SubscriberURL: "https://hapcky.me",
		TTL:           3600,
	}
}

type PushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url,omitempty"`
	Tag   string `json:"tag,omitempty"`
}

type Worker struct {
	DB     *sql.DB
	Config *Config
}

func NewWorker(db *sql.DB, config *Config) *Worker {
	if config == nil {
		config = DefaultConfig()
	}

	if config.VAPIDPublic == "" {
		config.VAPIDPublic = os.Getenv("VAPID_PUBLIC")
	}
	if config.VAPIDPrivate == "" {
		config.VAPIDPrivate = os.Getenv("VAPID_PRIVATE")
	}

	return &Worker{DB: db, Config: config}
}

func (w *Worker) Start(ctx context.Context) {
	log.Printf("[push-worker] 워커 시작됨 - %v마다 배치 처리", w.Config.BatchInterval)
	ticker := time.NewTicker(w.Config.BatchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := w.processBatch(ctx, w.Config.BatchSize); err != nil {
				log.Printf("[push-worker] 배치 처리 오류: %v", err)
			}
		case <-ctx.Done():
			log.Printf("[push-worker] 워커 종료됨")
			return
		}
	}
}

type dueJob struct {
	ID       int64
	MakeupID string
	UserID   string
	Kind     string
}

func (w *Worker) processBatch(ctx context.Context, limit int) error {
	jobs, err := w.fetchPendingJobs(ctx, limit)
	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}

	return w.processJobs(ctx, jobs)
}

func (w *Worker) fetchPendingJobs(ctx context.Context, limit int) ([]dueJob, error) {
	tx, err := w.DB.BeginTx(ctx, nil)
	if err != nil {
		log.Printf("[push-worker] 트랜잭션 시작 실패: %v", err)
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	rows, err := tx.QueryContext(ctx, `
	  WITH picked AS (
	    SELECT id
	    FROM scheduled_notifications
	    WHERE status='pending' AND scheduled_at <= now()
	    ORDER BY scheduled_at
	    FOR UPDATE SKIP LOCKED
	    LIMIT $1
	  )
	  UPDATE scheduled_notifications sn
	  SET status='processing'
	  FROM picked
	  WHERE sn.id = picked.id
	  RETURNING sn.id, sn.makeup_id, sn.user_id, sn.kind
	`, limit)
	if err != nil {
		log.Printf("[push-worker] 작업 조회 실패: %v", err)
		return nil, err
	}
	defer rows.Close()

	var jobs []dueJob
	for rows.Next() {
		var j dueJob
		if err := rows.Scan(&j.ID, &j.MakeupID, &j.UserID, &j.Kind); err != nil {
			log.Printf("[push-worker] 작업 스캔 실패: %v", err)
			return nil, err
		}
		jobs = append(jobs, j)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[push-worker] 트랜잭션 커밋 실패: %v", err)
		return nil, err
	}

	return jobs, nil
}

func (w *Worker) processJobs(ctx context.Context, jobs []dueJob) error {
	for _, j := range jobs {
		if err := w.sendMakeupNotification(ctx, j); err != nil {
			log.Printf("[push-worker] 작업 %d 처리 실패: %v", j.ID, err)
			w.failJob(ctx, j.ID, err)
		} else {
			w.successJob(ctx, j.ID)
		}
	}
	return nil
}

func (w *Worker) sendMakeupNotification(ctx context.Context, j dueJob) error {
	subscriptions, err := w.fetchUserSubscriptions(ctx, j.UserID)
	if err != nil {
		return err
	}

	payload := w.createNotificationPayload(j.Kind)
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return w.sendToAllSubscriptions(ctx, j.ID, subscriptions, data)
}

func (w *Worker) fetchUserSubscriptions(ctx context.Context, userID string) ([]*webpush.Subscription, error) {
	rows, err := w.DB.QueryContext(ctx, `
	  SELECT endpoint, p256dh, auth
	  FROM push_subscriptions
	  WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*webpush.Subscription
	for rows.Next() {
		var endpoint, p256dh, auth string
		if err := rows.Scan(&endpoint, &p256dh, &auth); err != nil {
			log.Printf("[push-worker] 구독 스캔 실패 - %v", err)
			continue
		}

		sub := &webpush.Subscription{
			Endpoint: endpoint,
			Keys:     webpush.Keys{P256dh: p256dh, Auth: auth},
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (w *Worker) createNotificationPayload(kind string) map[string]string {
	title, body, tag := payloadForKind(kind)
	return map[string]string{
		"title": title,
		"body":  body,
		"url":   "/calendar",
		"tag":   tag,
	}
}

func (w *Worker) sendToAllSubscriptions(ctx context.Context, jobID int64, subscriptions []*webpush.Subscription, data []byte) error {
	for _, sub := range subscriptions {
		if err := w.sendToSubscription(ctx, jobID, sub, data); err != nil {
			log.Printf("[push-worker] 작업 %d: 구독 %s 전송 실패 - %v", jobID, sub.Endpoint, err)
		}
	}
	return nil
}

func (w *Worker) sendToSubscription(ctx context.Context, jobID int64, sub *webpush.Subscription, data []byte) error {
	resp, err := webpush.SendNotification(data, sub, &webpush.Options{
		Subscriber:      w.Config.SubscriberURL,
		VAPIDPublicKey:  w.Config.VAPIDPublic,
		VAPIDPrivateKey: w.Config.VAPIDPrivate,
		TTL:             w.Config.TTL,
	})

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return w.handleSendError(ctx, jobID, sub.Endpoint, resp, err)
	}

	return nil
}

func (w *Worker) handleSendError(ctx context.Context, jobID int64, endpoint string, resp *http.Response, err error) error {
	if resp != nil {
		if resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound {
			log.Printf("[push-worker] 작업 %d: 만료된 구독 삭제 - %s", jobID, endpoint)
			_, _ = w.DB.ExecContext(ctx, "DELETE FROM push_subscriptions WHERE endpoint=$1", endpoint)
			return nil
		}

		if body, readErr := io.ReadAll(resp.Body); readErr == nil {
			log.Printf("[push-worker] 작업 %d: 푸시 전송 실패 - %v, %s", jobID, err, string(body))
		}
	}

	return err
}

func payloadForKind(kind string) (title, body, tag string) {
	nowKST := time.Now().In(mustKST())
	switch kind {
	case "makeup-1d":
		return "내일 보강 알림", "내일 예정된 보강이 있어요. 시간을 다시 한 번 확인해 주세요.", "makeup-1d"
	case "makeup-30m":
		return "보강 30분 전", "보강 시작까지 30분 남았어요. 늦지 않게 준비해주세요.", "makeup-30m"
	default:
		return "알림", nowKST.Format("15:04") + " 알림", "general"
	}
}

func mustKST() *time.Location {
	loc, _ := time.LoadLocation("Asia/Seoul")
	return loc
}

func (w *Worker) successJob(ctx context.Context, id int64) {
	_, err := w.DB.ExecContext(ctx, `
	  UPDATE scheduled_notifications
	  SET status='sent', attempts=attempts+1, sent_at=now(), last_error=NULL
	  WHERE id=$1
	`, id)
	if err != nil {
		log.Printf("[push-worker] 작업 %d: 성공 상태 업데이트 실패 - %v", id, err)
	}
}

func (w *Worker) failJob(ctx context.Context, id int64, e error) {
	log.Printf("[push-worker] 작업 %d: 실패 상태로 업데이트 - %v", id, e)
	_, err := w.DB.ExecContext(ctx, `
	  UPDATE scheduled_notifications
	  SET attempts = attempts + 1,
	      status = CASE WHEN attempts + 1 >= $2 THEN 'failed' ELSE 'pending' END,
	      last_error = $3
	  WHERE id=$1
	`, id, w.Config.MaxRetries, e.Error())
	if err != nil {
		log.Printf("[push-worker] 작업 %d: 실패 상태 업데이트 실패 - %v", id, err)
	}
}
