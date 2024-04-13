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
    parent_dir BIGINT REFERENCES dirs (dir_id) ON DELETE CASCADE DEFAULT NULL,
    icon_url VARCHAR(512) NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS notes
(
    note_id     BIGSERIAL PRIMARY KEY,
    name        VARCHAR(256) NOT NULL,
    body        TEXT         NOT NULL           DEFAULT '' NOT NULL,
    created_at  TIMESTAMP                       DEFAULT NOW() NOT NULL,
    last_update TIMESTAMP                       DEFAULT NOW() NOT NULL,
    parent_dir  BIGINT REFERENCES dirs (dir_id)  ON DELETE CASCADE DEFAULT NULL,
    repeated_num BIGINT NOT NULL DEFAULT 0,
    user_id     BIGINT REFERENCES users (user_id) NOT NULL
);


CREATE TABLE IF NOT EXISTS snippets
(
    snippet_id  BIGSERIAL PRIMARY KEY,
    name        VARCHAR(128) NOT NULL,
    description TEXT DEFAULT '' NOT NULL, 
    body        TEXT         NOT NULL DEFAULT '',
    user_id     BIGINT REFERENCES users (user_id)  ON DELETE CASCADE DEFAULT NULL 
);


-- TRIGGERS
create or replace function update_note_repeated_num()
returns trigger
as
    $$
BEGIN
    NEW.repeated_num = (SELECT COUNT(*) FROM notes WHERE name = NEW.name AND user_id = NEW.user_id AND parent_dir = NEW.parent_dir);
    RETURN NEW;
END;
$$
language plpgsql
;


CREATE TRIGGER update_note_repeated_num_trigger
BEFORE UPDATE ON notes
FOR EACH ROW
EXECUTE FUNCTION update_note_repeated_num();

create or replace function update_dir_repeated_num()
returns trigger
as
    $$
BEGIN
    NEW.repeated_num = (SELECT COUNT(*) FROM dirs WHERE name = NEW.name AND user_id = NEW.user_id AND parent_dir = NEW.parent_dir);
    RETURN NEW;
END;
$$
language plpgsql
;


CREATE TRIGGER update_dir_repeated_num_trigger
BEFORE UPDATE ON dirs
FOR EACH ROW
EXECUTE FUNCTION update_dir_repeated_num();
