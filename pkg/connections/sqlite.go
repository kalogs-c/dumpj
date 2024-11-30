package connections

import (
	"context"
	"database/sql"

	"github.com/kalogs-c/dumpj/internal/db"
	dumpjsql "github.com/kalogs-c/dumpj/sql"
	_ "github.com/mattn/go-sqlite3"
	goose "github.com/pressly/goose/v3"
)

type closeConnFunc func()

func migrate(conn *sql.DB) error {
	goose.SetBaseFS(dumpjsql.MigrationsEmbed)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(conn, "migrations"); err != nil {
		return err
	}

	return nil
}

func NewSQLite(ctx context.Context, path string) (*db.Queries, closeConnFunc, error) {
	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, nil, err
	}

	if err := migrate(conn); err != nil {
		return nil, nil, err
	}

	queries := db.New(conn)

	return queries, func() { conn.Close() }, nil
}
