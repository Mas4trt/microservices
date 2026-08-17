-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    order_uuid       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid        UUID NOT NULL,
    total_price      NUMERIC(14, 2) NOT NULL CHECK (total_price >= 0),
    transaction_uuid UUID NULL,
    payment_method   TEXT NOT NULL DEFAULT 'UNKNOWN'
        CHECK (payment_method IN ('UNKNOWN', 'CARD', 'SBP', 'CREDIT_CARD', 'INVESTOR_MONEY')),
    status           TEXT NOT NULL DEFAULT 'PENDING_PAYMENT'
        CHECK (status IN ('PENDING_PAYMENT', 'PAID', 'CANCELLED')),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

CREATE INDEX idx_orders_user_uuid ON orders (user_uuid);

CREATE INDEX idx_orders_status ON orders (status);

-- Автоматическое обновление updated_at при UPDATE
-- +goose StatementBegin
CREATE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_orders_set_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_orders_set_updated_at ON orders;

DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS orders;