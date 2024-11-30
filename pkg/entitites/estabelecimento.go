package entitites

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/kalogs-c/dumpj/internal/db"
)

type Estabelecimento struct {
	CnpjBasico                string         `csv_column:"1"`
	CnpjOrdem                 string         `csv_column:"2"`
	CnpjDv                    string         `csv_column:"3"`
	IdentificadorMatrizFilial int64          `csv_column:"4"`
	NomeFantasia              sql.NullString `csv_column:"5"`
	Situacao                  string         `csv_column:"6"`
	DataAbertura              time.Time      `csv_column:"11"`
	Cnae                      string         `csv_column:"12"`
	TipoLogradouro            sql.NullString `csv_column:"14"`
	Logradouro                string         `csv_column:"15"`
	Numero                    string         `csv_column:"16"`
	Complemento               sql.NullString `csv_column:"17"`
	Bairro                    string         `csv_column:"18"`
	Cep                       string         `csv_column:"19"`
	Uf                        string         `csv_column:"20"`
	Municipio                 string         `csv_column:"21"`
	Ddd                       sql.NullString `csv_column:"22"`
	Telefone                  sql.NullString `csv_column:"23"`
	Email                     sql.NullString `csv_column:"28"`
}

func (e *Estabelecimento) Save(ctx context.Context, q *db.Queries) error {
	attr := db.CreateEstabelecimentoParams{
		CnpjBasico:                e.CnpjBasico,
		CnpjOrdem:                 e.CnpjOrdem,
		CnpjDv:                    e.CnpjDv,
		IdentificadorMatrizFilial: e.IdentificadorMatrizFilial,
		NomeFantasia:              e.NomeFantasia,
		DataAbertura:              e.DataAbertura,
		Cnae:                      e.Cnae,
		Logradouro:                e.Logradouro,
		Numero:                    e.Numero,
		Complemento:               e.Complemento,
		Bairro:                    e.Bairro,
		Cep:                       e.Cep,
		Uf:                        e.Uf,
		Municipio:                 e.Municipio,
		Ddd:                       e.Ddd,
		Telefone:                  e.Telefone,
		Email:                     e.Email,
	}

	validMap, ok := ctx.Value("validEstabelecimentos").(map[string]bool)
	if ok {
		validMap[e.CnpjBasico] = true
	}

	return q.CreateEstabelecimento(ctx, attr)
}

func (e *Estabelecimento) IsValid() bool {
	ddd := e.Ddd.String
	telefone := e.Telefone.String

	ignore := []string{"", "0", "00", "**"}
	for _, v := range ignore {
		if v == ddd || v == telefone {
			return false
		}
	}

	return e.Situacao != "02"
}

func (e *Estabelecimento) Print() {
	fmt.Printf("Estabelecimento: %s-%s.%s.\n", e.CnpjBasico, e.CnpjOrdem, e.CnpjDv)
}
