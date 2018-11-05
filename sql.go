package snippets

import (
	"context"
	"database/sql"
)

func query(ctx context.Context, db *sql.DB, id int, limit int) ([]int64, error) {
	rows, err := db.QueryContext(ctx, "SELECT id FROM users WHERE id > ? LIMIT ?", id, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	results := make([]int64, 0, limit)
	for rows.Next() {
		var id int64
		rows.Scan(&id)
		results = append(results, id)
	}
	return results, nil
}

func queryRow(ctx context.Context, db *sql.DB, id int) (string, error) {
	var name string
	err := db.QueryRowContext(ctx, "SELECT name FROM users WHERE id = ?").Scan(&name)
	return name, err
}
