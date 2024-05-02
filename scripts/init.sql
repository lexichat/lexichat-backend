CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    user_id VARCHAR(25) UNIQUE,
    user_name VARCHAR(25),
    phone_number VARCHAR(15),
    fcm_token TEXT,
    profile_picture BYTEA,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_userid ON users (user_id);



CREATE TABLE IF NOT EXISTS channels (
    id BIGSERIAL PRIMARY KEY,
    channel_name VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    tonality_tag VARCHAR,
    description TEXT
);



CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    channel_id BIGSERIAL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sender_user_id VARCHAR,
    message TEXT,
    status VARCHAR,
    FOREIGN KEY (channel_id) REFERENCES channels(id),
    FOREIGN KEY (sender_user_id) REFERENCES users(user_id)
);

CREATE INDEX IF NOT EXISTS idx_channel_id_created_at ON messages (channel_id, created_at DESC);


CREATE TABLE IF NOT EXISTS channel_users (
    channel_id BIGINT,
    user_id VARCHAR,
    FOREIGN KEY (channel_id) REFERENCES channels(id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    PRIMARY KEY (channel_id, user_id)
);