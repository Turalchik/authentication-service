CREATE TABLE sessions (
    user_id UUID PRIMARY KEY,
    refresh_token TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    ip_addr TEXT NOT NULL
);