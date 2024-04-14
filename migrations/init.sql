CREATE TABLE IF NOT EXISTS banner(
    banner_id  SERIAL PRIMARY KEY,
    content    BYTEA NOT NULL,
    is_active  BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS tag(
    tag_id  INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS feature(
    feature_id  INT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS banner_tag_feature(
    banner_id INT,
    tag_id INT,
    feature_id INT,
    FOREIGN KEY (banner_id) REFERENCES banner(banner_id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tag(tag_id) ON DELETE CASCADE,
    FOREIGN KEY (feature_id) REFERENCES feature(feature_id) ON DELETE CASCADE,
    CONSTRAINT PK_TagFeature PRIMARY KEY (tag_id, feature_id)
);

CREATE TABLE IF NOT EXISTS "user"(
    user_id    SERIAL PRIMARY KEY,
    login      VARCHAR(32) UNIQUE  NOT NULL,
    password   VARCHAR(32) NOT NULL,
    is_admin   BOOLEAN DEFAULT FALSE,
    tag_id     INT,
    FOREIGN KEY (tag_id) REFERENCES tag(tag_id) ON DELETE SET NULL
);


CREATE INDEX index_banner
ON banner(banner_id);


CREATE INDEX index_feature_tag
ON banner_tag_feature(tag_id, feature_id);

INSERT INTO banner (
    content, is_active
)
SELECT
    ('{"url": "u://banner/' || i::text || '", "title": "title of banner ' ||  i::text  || '"}')::bytea,
    true
FROM generate_series(1, 1000000) s(i);

INSERT INTO tag (
    tag_id
)
SELECT
    i
FROM generate_series(1, 10000) s(i);

INSERT INTO feature (
    feature_id
)
SELECT
    i
FROM generate_series(1, 10000) s(i);

INSERT INTO banner_tag_feature(
    banner_id, tag_id, feature_id
)
SELECT
    1000 * (tag.num - 1) + feature.num,
    tag.num,
    feature.num
FROM
    generate_series(1, 1000) AS tag(num),
    generate_series(1, 1000) AS feature(num);

--создать еще индексы
