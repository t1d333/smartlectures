package pg

const (
	InsertNewDirCMD = `
		INSERT INTO dirs(name, user_id, parent_dir)
		VALUES($1, $2, CASE WHEN $3=0 THEN NULL ELSE $3 END)
		RETURNING dir_id;
	`

	SelectDirByIDCMD = `
		SELECT dir_id, name, user_id, parent_dir
		FROM dirs
		WHERE dir_id = $1;
	`

	UpdateDirCMD= `
		UPDATE dirs
		SET name = $2, parent_dir = CASE WHEN $3=0 THEN NULL ELSE $3 END
		WHERE dir_id = $1
		RETURNING dir_id;
	`

	DeleteDirCMD = `
		DELETE FROM dirs
		WHERE dir_id = $1;
	`

	SelectUserDirsOverview = `
		SELECT dir_id
		FROM dirs
		WHERE user_id = $1 AND parent_dir IS NULL;
	`

	SelectSubdirs = `
		SELECT dir_id 
		FROM dirs
		WHERE parent_dir = $1;
	`
)
