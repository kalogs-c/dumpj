// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package db

import (
	"database/sql"
	"time"
)

type Cnae struct {
	Codigo    string
	Descricao string
}

type Ddd struct {
	ID int64
	Uf string
}

type Empresa struct {
	CnpjBasico       string
	RazaoSocial      string
	NaturezaJuridica string
	CapitalSocial    int64
	PorteEmpresa     string
}

type Estabelecimento struct {
	ID                        int64
	CnpjBasico                string
	CnpjOrdem                 string
	CnpjDv                    string
	IdentificadorMatrizFilial int64
	NomeFantasia              sql.NullString
	DataAbertura              time.Time
	Cnae                      string
	Logradouro                string
	Numero                    string
	Complemento               sql.NullString
	Bairro                    string
	Cep                       string
	Uf                        string
	Municipio                 string
	Ddd                       sql.NullString
	Telefone                  sql.NullString
	Email                     sql.NullString
}

type Municipio struct {
	Codigo    string
	Descricao string
}

type NaturezasJuridica struct {
	Codigo    string
	Descricao string
}

type Uf struct {
	Codigo    string
	Descricao string
}
