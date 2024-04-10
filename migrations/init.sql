CREATE TABLE IF NOT EXISTS banner_tag_feature(
    banner_id INT,
    tag_id INT,
    feature_id INT,
    FOREIGN KEY (banner_id) REFERENCES banner(banner_id) ON DELETE CASCADE,
    CONSTRAINT PK_TagFeature PRIMARY KEY (tag_id, feature_id)
);

CREATE TABLE IF NOT EXISTS banner(
    banner_id  SERIAL PRIMARY KEY,
    content    BYTEA NOT NULL,
    is_active  BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "user"(
    user_id    SERIAL PRIMARY KEY,
    login      VARCHAR(32) NOT NULL,
    password   VARCHAR(32) NOT NULL,
    is_admin   BOOLEAN DEFAULT FALSE,
    tag_id     INT
);
