-- +goose Up
-- +goose StatementBegin
CREATE TABLE municipios (
  codigo    text    PRIMARY KEY,
  descricao text    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE municipios;
-- +goose StatementEnd
