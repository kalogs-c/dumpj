package entitites

import (
	"context"
	"fmt"

	"github.com/kalogs-c/dumpj/internal/db"
)

type Municipio struct {
	Codigo    string `csv_column:"1"`
	Descricao string `csv_column:"2"`
}

func (m *Municipio) Save(ctx context.Context, q *db.Queries) error {
	return q.CreateMunicipio(ctx, db.CreateMunicipioParams(*m))
}

func (m *Municipio) IsValid() bool {
	return m.Codigo != "" && m.Descricao != ""
}

func (m *Municipio) Print() {
	fmt.Printf("Municipio: %s\n", m.Descricao)
}
