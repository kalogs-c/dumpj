package connections

import (
	"context"
	"database/sql"

	"github.com/kalogs-c/dumpj/internal/sql/db"
)

func NewSQLite(ctx context.Context, path string) (*db.Queries, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	queries := db.New(conn)

	return queries, nil
}
