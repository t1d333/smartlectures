package pg

const (
	SelectUserSnippets = `
		SELECT snippet_id, name, description, body, user_id 
		FROM snippets
		WHERE user_id IS NULL OR user_id = $1
	`
)
