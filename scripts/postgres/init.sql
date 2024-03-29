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
    user_id    BIGINT REFERENCES users (user_id) NOT NULL,
    repeated_num BIGINT DEFAULT 0,
    parent_dir BIGINT REFERENCES dirs (dir_id) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS notes
(
    note_id     BIGSERIAL PRIMARY KEY,
    name        VARCHAR(256) NOT NULL,
    body        TEXT         NOT NULL           DEFAULT '' NOT NULL,
    created_at  TIMESTAMP                       DEFAULT NOW() NOT NULL,
    last_update TIMESTAMP                       DEFAULT NOW() NOT NULL,
    parent_dir  BIGINT REFERENCES dirs (dir_id) DEFAULT NULL,
    repeated_num BIGINT NOT NULL DEFAULT 0,
    user_id     BIGINT REFERENCES users (user_id) NOT NULL

);


-- TRIGGERS

CREATE OR REPLACE FUNCTION update_note_repeated_num()
RETURNS TRIGGER AS $$
BEGIN
    NEW.repeated_num = (SELECT COUNT(*) FROM notes WHERE name = NEW.name AND user_id = NEW.user_id AND parent_dir = NEW.parent_dir);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_note_repeated_num_trigger
BEFORE UPDATE ON notes
FOR EACH ROW
EXECUTE FUNCTION update_note_repeated_num();

CREATE OR REPLACE FUNCTION update_dir_repeated_num()
RETURNS TRIGGER AS $$
BEGIN
    NEW.repeated_num = (SELECT COUNT(*) FROM dirs WHERE name = NEW.name AND user_id = NEW.user_id AND parent_dir = NEW.parent_dir);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_dir_repeated_num_trigger
BEFORE UPDATE ON dirs
FOR EACH ROW
EXECUTE FUNCTION update_dir_repeated_num();

