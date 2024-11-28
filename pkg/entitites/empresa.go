package entitites

import (
	"context"
	"fmt"

	"github.com/kalogs-c/dumpj/internal/db"
)

type Empresa struct {
	CnpjBasico       string `csv_column:"1"`
	RazaoSocial      string `csv_column:"2"`
	NaturezaJuridica string `csv_column:"3"`
	CapitalSocial    int64  `csv_column:"5"`
	PorteEmpresa     string `csv_column:"6"`
}

func (e *Empresa) Save(ctx context.Context, q *db.Queries) error {
	return q.CreateEmpresa(ctx, db.CreateEmpresaParams(*e))
}

func (e *Empresa) IsValid() bool {
	return e.CnpjBasico != ""
}

func (e *Empresa) Print() {
	fmt.Printf("Empresa: %s\n", e.CnpjBasico)
}
