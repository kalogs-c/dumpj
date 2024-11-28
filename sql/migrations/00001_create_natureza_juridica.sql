-- +goose Up
-- +goose StatementBegin
CREATE TABLE naturezas_juridicas (
  codigo    text    PRIMARY KEY,
  descricao text    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE naturezas_juridicas;
-- +goose StatementEnd
