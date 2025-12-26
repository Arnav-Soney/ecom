-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products(
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price > 0),
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
