-- +goose Up
-- +goose StatementBegin
CREATE TABLE ddd (
  id INTEGER PRIMARY KEY,
  uf TEXT NOT NULL,
  -- FOREIGN KEY (uf) REFERENCES ufs (codigo)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ddd;
-- +goose StatementEnd
