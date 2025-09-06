package services

import (
	"errors"
	"julia/internal/models"
	"julia/internal/repositories"
	"julia/utils"
	"net/http"
	"sync"
	"time"
)

type Config struct {
	SessionTTL time.Duration
	HMACSecret []byte
	CookieName string
	CookiePath string
	Secure     bool
	SameSite   http.SameSite
	Domain     string
}

type Session struct {
	UserID    string
	ExpiresAt time.Time
}

type AuthService interface {
	CreateSession(userID string) (string, *Session, error)
	GetSession(sid string) (*Session, bool)
	DeleteSession(sid string)
	ExtendSession(sid string) (*Session, bool)

	IssueSessionCookie(w http.ResponseWriter, sid string)
	ClearSessionCookie(w http.ResponseWriter)

	ParseAndVerifySID(rawCookie string) (sid string, err error)

	VerifyCredential(id, password string) (string, error)
	GetUserByID(userID string) (*models.User, error)

	SessionName() string
}

type authService struct {
	userRepo repositories.UserRepository
	cfg      Config

	sessionStore map[string]*Session
	mu           sync.RWMutex
}

func NewAuthService(userRepo repositories.UserRepository, cfg Config) *authService {
	return &authService{
		userRepo:     userRepo,
		cfg:          cfg,
		sessionStore: make(map[string]*Session),
	}
}

func (s *authService) SessionName() string {
	return s.cfg.CookieName
}

func (s *authService) CreateSession(userID string) (string, *Session, error) {
	sid, err := utils.NewSID()
	if err != nil {
		return "", nil, err
	}
	session := &Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.cfg.SessionTTL),
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessionStore[sid] = session
	return sid, session, nil
}

func (s *authService) GetSession(sid string) (*Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessionStore[sid]
	if !ok {
		return nil, false
	}
	if time.Now().After(session.ExpiresAt) {
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.sessionStore, sid)
		return nil, false
	}
	return session, true
}

func (s *authService) DeleteSession(sid string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessionStore, sid)
}

func (s *authService) ExtendSession(sid string) (*Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessionStore[sid]
	if !ok {
		return nil, false
	}
	session.ExpiresAt = time.Now().Add(s.cfg.SessionTTL)
	return session, true
}

func (s *authService) IssueSessionCookie(w http.ResponseWriter, sid string) {
	utils.SetCookie(w, utils.EncodeCookie(sid, s.cfg.HMACSecret), utils.Cookie{
		Name:     s.cfg.CookieName,
		Path:     s.cfg.CookiePath,
		MaxAge:   int(s.cfg.SessionTTL.Seconds()),
		HttpOnly: false, // 개발 환경에서는 false로 설정
		Secure:   s.cfg.Secure,
		SameSite: s.cfg.SameSite,
		Domain:   s.cfg.Domain,
	})
}

func (s *authService) ClearSessionCookie(w http.ResponseWriter) {
	utils.ClearCookie(w, utils.Cookie{
		Name:     s.cfg.CookieName,
		Path:     s.cfg.CookiePath,
		MaxAge:   0,
		HttpOnly: false, // 개발 환경에서는 false로 설정
		Secure:   s.cfg.Secure,
		SameSite: s.cfg.SameSite,
		Domain:   s.cfg.Domain,
	})
}

func (s *authService) ParseAndVerifySID(rawCookie string) (sid string, err error) {
	sid, sig, err := utils.DecodeCookie(rawCookie)
	if err != nil {
		return "", err
	}
	if !utils.VerifySID(sid, sig, s.cfg.HMACSecret) {
		return "", errors.New("invalid signature")
	}
	return sid, nil
}

func (s *authService) VerifyCredential(id, password string) (string, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return "", err
	}
	if !utils.ComparePassword(user.Password, password) {
		return "", errors.New("invalid password")
	}
	return user.ID, nil
}

func (s *authService) GetUserByID(userID string) (*models.User, error) {
	return s.userRepo.GetUserByID(userID)
}
