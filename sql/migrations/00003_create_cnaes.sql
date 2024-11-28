-- +goose Up
-- +goose StatementBegin
CREATE TABLE cnaes (
  codigo    text    PRIMARY KEY,
  descricao text    NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cnaes;
-- +goose StatementEnd
