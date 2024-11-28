package entitites

import (
	"context"
	"fmt"

	"github.com/kalogs-c/dumpj/internal/db"
)

type CNAE struct {
	Codigo    string `csv_column:"1"`
	Descricao string `csv_column:"2"`
}

func (c *CNAE) Save(ctx context.Context, q *db.Queries) error {
	return q.CreateCnae(ctx, db.CreateCnaeParams(*c))
}

func (c *CNAE) IsValid() bool {
	return c.Codigo != "" && c.Descricao != ""
}

func (c *CNAE) Print() {
	fmt.Printf("CNAE: %s\n", c.Codigo)
}
