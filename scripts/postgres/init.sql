CREATE TABLE IF NOT EXISTS users
(
    user_id  BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE  NOT NULL,
    email    VARCHAR(300) UNIQUE NOT NULL,
    password BYTEA               NOT NULL,
    name     VARCHAR(128) DEFAULT '',
    surname  VARCHAR(128) DEFAULT ''
);

CREATE TABLE IF NOT EXISTS dirs
(
    dir_id     BIGSERIAL PRIMARY KEY,
    name       VARCHAR(128) NOT NULL,
    user_id    BIGINT REFERENCES users (user_id),
    parent_dir BIGINT REFERENCES dirs (dir_id)
);

CREATE TABLE IF NOT EXISTS notes
(
    note_id     BIGSERIAL PRIMARY KEY,
    name        VARCHAR(256) NOT NULL,
    body        TEXT         NOT NULL           DEFAULT '',
    created_at  TIMESTAMP                       DEFAULT NOW(),
    last_update TIMESTAMP                       DEFAULT NULL,
    parent_dir  BIGINT REFERENCES dirs (dir_id) DEFAULT 0,
    user_id     BIGINT REFERENCES users (user_id)

);
