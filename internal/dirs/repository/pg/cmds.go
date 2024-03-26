package pg

const (
	InsertNewDirCMD = `
		INSERT INTO dirs(name, user_id, parent_dir, repeated_num)
		VALUES($1, $2, $3, (SELECT COUNT(*) FROM dirs WHERE name = $4 AND parent_dir = $5))
		RETURNING dir_id;
	`
	
	InsertNewDirWithNullParentCMD = `
		INSERT INTO dirs(name, user_id, repeated_num)
		VALUES($1, $2, (SELECT COUNT(*) FROM dirs WHERE name = $3 AND parent_dir IS NULL))
		RETURNING dir_id;
	`
	
	SelectDirByIDCMD = `
		SELECT dir_id, name, user_id, parent_dir, repeated_num
		FROM dirs
		WHERE dir_id = $1;
	`

	UpdateDirCMD = `
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
