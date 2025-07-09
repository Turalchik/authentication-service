CREATE TABLE sessions (
    user_id UUID PRIMARY KEY,
    refresh_token_hash TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    ip_addr TEXT NOT NULL
);