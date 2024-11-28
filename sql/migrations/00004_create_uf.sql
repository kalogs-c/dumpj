-- +goose Up
-- +goose StatementBegin
CREATE TABLE uf (
  codigo TEXT PRIMARY KEY,
  descricao TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE uf;
-- +goose StatementEnd
