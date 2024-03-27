package pg

const (
	InsertNewNoteCMD = `
		INSERT INTO notes(name, body, parent_dir, user_id, repeated_num)
		VALUES($1, $2, $3, $4, (SELECT COUNT(*) FROM notes WHERE name = $5 AND parent_dir = $6 AND user_id = $7))
		RETURNING note_id;
	`

	InsertNewNoteWithNullParentCMD = `
		INSERT INTO notes(name, body, user_id, repeated_num)
		VALUES($1, $2, $3, (SELECT COUNT(*) FROM notes WHERE name = $4 AND user_id = $5 AND parent_dir IS NULL))
		RETURNING note_id;
	`

	SelectNoteByIDCMD = `
		SELECT note_id, name, body, created_at, last_update, parent_dir, user_id, repeated_num
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
		SELECT note_id, name, parent_dir, repeated_num
		FROM notes
		WHERE user_id = $1;
	`
)
