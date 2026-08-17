package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, order model.OrderDto) (uuid.UUID, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate order uuid: %w", err)
	}
	order.OrderUUID = newUUID

	repoOrder := repoConverter.ServiceToRepoModel(order)
	orderParts := repoConverter.ServiceToRepoOrderParts(newUUID, order.PartUuids)

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, `
		INSERT INTO orders (order_uuid, user_uuid, total_price, transaction_uuid, payment_method, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`,
		repoOrder.OrderUUID, repoOrder.UserUUID, repoOrder.TotalPrice,
		repoOrder.TransactionUUID, repoOrder.PaymentMethod, repoOrder.Status,
	)
	if err != nil {
		return uuid.Nil, fmt.Errorf("insert order: %w", err)
	}

	if err := insertOrderParts(ctx, tx, orderParts); err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("commit tx: %w", err)
	}

	return newUUID, nil
}
