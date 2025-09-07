CREATE TYPE assignment_status AS ENUM ('pending', 'completed');

CREATE TABLE assignments (
    assignment_id SERIAL PRIMARY KEY,
    user_id       TEXT NOT NULL,
    content       TEXT NOT NULL,
    status        assignment_status NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);