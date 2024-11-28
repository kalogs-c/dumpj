package entitites

import (
	"context"
	"fmt"

	"github.com/kalogs-c/dumpj/internal/db"
)

type NaturezaJuridica struct {
	Descricao string `csv_column:"2"`
	Codigo    string `csv_column:"1"`
}

func (nj *NaturezaJuridica) Save(ctx context.Context, q *db.Queries) error {
	return q.CreateNaturezaJuridica(ctx, db.CreateNaturezaJuridicaParams(*nj))
}

func (nj *NaturezaJuridica) IsValid() bool {
	return nj.Codigo != "" && nj.Descricao != ""
}

func (nj *NaturezaJuridica) Print() {
	fmt.Printf("Natureza Juridica: %s\n", nj.Codigo)
}
