CREATE TABLE IF NOT EXISTS banner_tag_feature(
    banner_id INT,
    tag_id INT,
    feature_id INT,
    FOREIGN KEY (banner_id) REFERENCES banner(banner_id) ON DELETE CASCADE,
    CONSTRAINT PK_TagFeature PRIMARY KEY (tag_id, feature_id)
);

CREATE TABLE IF NOT EXISTS banner(
    banner_id  INT PRIMARY KEY,
    content    BYTEA,
    is_active  BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
