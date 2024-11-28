package connections

import (
	"context"
	"database/sql"

	"github.com/kalogs-c/dumpj/internal/db"
	_ "github.com/mattn/go-sqlite3"
)

type closeConnFunc func()

func NewSQLite(ctx context.Context, path string) (*db.Queries, closeConnFunc, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, nil, err
	}

	queries := db.New(conn)

	return queries, func() { conn.Close() }, nil
}
