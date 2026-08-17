package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	repoModel "github.com/Mas4trt/microservices/order/internal/repository/model"
)

// insertOrderParts вставляет строки order_parts батчем в рамках переданной транзакции.
func insertOrderParts(ctx context.Context, tx pgx.Tx, parts []repoModel.OrderPart) error {
	if len(parts) == 0 {
		return nil
	}

	batch := &pgx.Batch{}
	for _, p := range parts {
		batch.Queue(
			`INSERT INTO order_parts (order_uuid, part_uuid) VALUES ($1, $2)`,
			p.OrderUUID, p.PartUUID,
		)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range parts {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("insert order_parts: %w", err)
		}
	}

	return nil
}
