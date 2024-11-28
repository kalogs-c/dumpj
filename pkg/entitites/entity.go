package entitites

import (
	"context"
	"regexp"
	"strings"

	"github.com/kalogs-c/dumpj/internal/db"
	"github.com/kalogs-c/dumpj/pkg/filemanager"
)

type Entity interface {
	Save(context.Context, *db.Queries) error
	IsValid() bool
	Print()
}

func PickEntity(filename string) Entity {
	reg := regexp.MustCompile("Empresas|Estabelecimentos|Socios")
	if reg.MatchString(filename) {
		filename = filename[:len(filename)-1]
	}

	switch strings.ToLower(filename) {
	case "cnaes":
		return &CNAE{}
	case "empresas":
		return &Empresa{}
	case "estabelecimentos":
		return &Estabelecimento{}
	case "municipios":
		return &Municipio{}
	case "naturezas":
		return &NaturezaJuridica{}
	default:
		return nil
	}
}

func NewEntityFromCSV(t any, row []string) error {
	err := filemanager.BindFields(row, t)
	if err != nil {
		return err
	}

	return nil
}
