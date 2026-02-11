-- +goose Up
ALTER TABLE orders ADD COLUMN IF NOT EXISTS stripe_session_id VARCHAR(255);
CREATE INDEX IF NOT EXISTS idx_orders_stripe_session_id ON orders (stripe_session_id) WHERE stripe_session_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS idx_orders_stripe_session_id;
ALTER TABLE orders DROP COLUMN IF EXISTS stripe_session_id;
