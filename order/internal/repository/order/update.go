package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, orderUUID uuid.UUID, newOrder model.OrderDto) error {
	repoOrder := repoConverter.ServiceToRepoModel(newOrder)
	orderParts := repoConverter.ServiceToRepoOrderParts(orderUUID, newOrder.PartUuids)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	tag, err := tx.Exec(ctx, `
		UPDATE orders
		SET user_uuid = $2, total_price = $3, transaction_uuid = $4, payment_method = $5, status = $6
		WHERE order_uuid = $1
	`,
		orderUUID, repoOrder.UserUUID, repoOrder.TotalPrice,
		repoOrder.TransactionUUID, repoOrder.PaymentMethod, repoOrder.Status,
	)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return model.ErrOrderNotFound
	}

	if _, err := tx.Exec(ctx, `DELETE FROM order_parts WHERE order_uuid = $1`, orderUUID); err != nil {
		return fmt.Errorf("delete old order_parts: %w", err)
	}

	if err := insertOrderParts(ctx, tx, orderParts); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}
