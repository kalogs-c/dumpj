package sql

import (
	"embed"
)

//go:embed migrations/*.sql
var MigrationsEmbed embed.FS
