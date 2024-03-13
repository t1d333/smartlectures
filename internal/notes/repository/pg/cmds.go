package pg

//     note_id     BIGSERIAL PRIMARY KEY,
//     name        VARCHAR(256) NOT NULL,
//     body        TEXT         NOT NULL           DEFAULT '',
//     created_at  TIMESTAMP                       DEFAULT NOW(),
//     last_update TIMESTAMP                       DEFAULT NULL,
//     parent_dir  BIGINT REFERENCES dirs (dir_id) DEFAULT 0,
//     user_id     BIGINT REFERENCES users (user_id)
//

const (
	InsertNewNoteCMD = `
		INSERT INTO notes(name, body, parent_dir, user_id)
		VALUES($1, $2, $3);
	`

	SelectNoteByIDCMD = `
		SELECT note_id, name, body, created_at, last_update, parent_dir, user_id
		FROM notes
		WHERE note_id = $1
	`

	UpdateNoteCMD = `
		UPDATE notes
		SET name = $2, body = $3, parent_dir = $4, last_update = NOW()
		WHERE note_id = $1
	`

	DeleteNodeCMD = `
		DELETE FROM notes
		WHERE note_id = $1
	`

	GetUserNotes = `
		SELECT note_id, name, parent_dir
		FROM notes
		WHERE user_id = $1
	`
)
