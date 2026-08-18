package order

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"

	repoModel "github.com/Mas4trt/microservices/order/internal/repository/model"
)

// insertOrderParts вставляет строки order_parts батчем в рамках переданной транзакции.
func insertOrderParts(ctx context.Context, tx pgx.Tx, parts []repoModel.OrderPart) (err error) {
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
	defer func() {
		if closeErr := br.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("close batch results: %w", closeErr)
		}
	}()

	for range parts {
		if _, err := br.Exec(); err != nil {
			return fmt.Errorf("insert order_parts: %w", err)
		}
	}

	return nil
}

func rollbackTx(ctx context.Context, tx pgx.Tx) {
	if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
		log.Printf("failed to rollback transaction: %v", err)
	}
}
