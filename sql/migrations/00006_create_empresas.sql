-- +goose Up
-- +goose StatementBegin
CREATE TABLE empresas (
  cnpj_basico TEXT PRIMARY KEY,
  razao_social TEXT NOT NULL,
  natureza_juridica TEXT NOT NULL,
  capital_social INTEGER DEFAULT 0 NOT NULL,
  porte_empresa TEXT NOT NULL,

  -- FOREIGN KEY (natureza_juridica) REFERENCES naturezas_juridicas (codigo)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE empresas;
-- +goose StatementEnd
