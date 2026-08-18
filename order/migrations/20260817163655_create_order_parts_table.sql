-- +goose Up
-- +goose StatementBegin
CREATE TABLE order_parts (
    order_uuid UUID NOT NULL REFERENCES orders (order_uuid) ON DELETE CASCADE,
    part_uuid  UUID NOT NULL,
    PRIMARY KEY (order_uuid, part_uuid)
);
-- +goose StatementEnd

CREATE INDEX idx_order_parts_part_uuid ON order_parts (part_uuid);

-- +goose Down
DROP TABLE IF EXISTS order_parts;