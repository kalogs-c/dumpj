-- +goose Up
-- +goose StatementBegin
CREATE TABLE estabelecimentos (
  id INTEGER PRIMARY KEY,
  cnpj_basico TEXT NOT NULL,
  cnpj_ordem TEXT NOT NULL,
  cnpj_dv TEXT NOT NULL,
  identificador_matriz_filial INTEGER CHECK (identificador_matriz_filial IN (1, 2)) NOT NULL, -- 1: Matriz, 2: Filial
  nome_fantasia TEXT,
  data_abertura DATE NOT NULL,
  cnae TEXT NOT NULL,
  tipo_logradouro TEXT,
  logradouro TEXT NOT NULL,
  numero TEXT NOT NULL,
  complemento TEXT,
  bairro TEXT NOT NULL,
  cep TEXT NOT NULL,
  uf TEXT NOT NULL,
  municipio TEXT NOT NULL,
  ddd TEXT,
  telefone TEXT,
  email TEXT
  -- FOREIGN KEY (cnpj_basico) REFERENCES empresas (cnpj_basico),
  -- FOREIGN KEY (cnae) REFERENCES cnaes (codigo),
  -- FOREIGN KEY (municipio) REFERENCES municipios (codigo),
  -- FOREIGN KEY (uf) REFERENCES ufs (codigo)
);

CREATE INDEX idx_estabelecimentos_cnpj ON estabelecimentos (cnpj_basico);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE estabelecimentos;
-- +goose StatementEnd
