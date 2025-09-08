CREATE TYPE assignment_status AS ENUM ('pending', 'completed');

CREATE TABLE assignments (
    assignment_id SERIAL PRIMARY KEY,
    user_id       TEXT NOT NULL,
    makeup_id     BIGINT NULL,
    content       TEXT NOT NULL,
    status        assignment_status NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_makeup
        FOREIGN KEY (makeup_id)
        REFERENCES makeups(makeup_id)
        ON DELETE CASCADE
);

ALTER TABLE assignments
  ALTER COLUMN created_at TYPE timestamptz
  USING created_at AT TIME ZONE 'Asia/Seoul',
  ALTER COLUMN updated_at TYPE timestamptz
  USING updated_at AT TIME ZONE 'Asia/Seoul';