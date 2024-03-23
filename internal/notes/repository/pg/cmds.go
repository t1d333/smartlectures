package pg

const (
	InsertNewNoteCMD = `
		INSERT INTO notes(name, body, parent_dir, user_id)
		VALUES($1, $2, CASE WHEN $3=0 THEN NULL ELSE $3 END, $4)
		RETURNING note_id;
	`

	SelectNoteByIDCMD = `
		SELECT note_id, name, body, created_at, last_update, parent_dir, user_id
		FROM notes
		WHERE note_id = $1;
	`

	UpdateNoteCMD = `
		UPDATE notes
		SET name = $2, body = $3, parent_dir = CASE WHEN $4=0 THEN NULL ELSE $4 END, last_update = NOW()
		WHERE note_id = $1
		RETURNING note_id;
	`

	DeleteNodeCMD = `
		DELETE FROM notes
		WHERE note_id = $1;
	`

	SelectUserNotesOverview = `
		SELECT note_id, name, parent_dir
		FROM notes
		WHERE user_id = $1;
	`
)
