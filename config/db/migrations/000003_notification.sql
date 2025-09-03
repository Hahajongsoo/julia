CREATE TABLE IF NOT EXISTS scheduled_notifications (
  id               BIGSERIAL PRIMARY KEY,
  makeup_id        BIGINT      NOT NULL,
  user_id          TEXT      NOT NULL,
  kind             TEXT      NOT NULL CHECK (kind IN ('makeup-1d','makeup-30m')),
  scheduled_at     TIMESTAMPTZ NOT NULL,
  status           TEXT      NOT NULL DEFAULT 'pending' CHECK (status IN ('pending','processing','sent','failed')),
  attempts         INT       NOT NULL DEFAULT 0,
  last_error       TEXT,
  sent_at          TIMESTAMPTZ,
  created_at       TIMESTAMPTZ DEFAULT now(),
  UNIQUE (makeup_id, kind),
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (makeup_id) REFERENCES makeups(makeup_id)
);