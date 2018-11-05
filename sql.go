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

func upsertRows(ctx context.Context, db *sql.DB, rows []string) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO users (name) VALUES (?) ON DUPLICATE KEY UPDATE name = ?")
	if err != nil {
		return err
	}

	for i := 0; i < len(rows); i++ {
		_, err := stmt.ExecContext(
			ctx,
			rows[i],
			rows[i],
		)
		if err != nil {
			if errR := tx.Rollback(); errR != nil {
				return errR
			}
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
