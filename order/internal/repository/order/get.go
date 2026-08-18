package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
	repoModel "github.com/Mas4trt/microservices/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.OrderDto, error) {
	var repoOrder repoModel.OrderDto
	err := r.pool.QueryRow(ctx, `
		SELECT order_uuid, user_uuid, total_price, transaction_uuid, payment_method, status
		FROM orders WHERE order_uuid = $1
	`, orderUUID).Scan(
		&repoOrder.OrderUUID, &repoOrder.UserUUID, &repoOrder.TotalPrice,
		&repoOrder.TransactionUUID, &repoOrder.PaymentMethod, &repoOrder.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.OrderDto{}, model.ErrOrderNotFound
		}
		return model.OrderDto{}, fmt.Errorf("query order: %w", err)
	}

	rows, err := r.pool.Query(ctx, `SELECT part_uuid FROM order_parts WHERE order_uuid = $1`, orderUUID)
	if err != nil {
		return model.OrderDto{}, fmt.Errorf("query order_parts: %w", err)
	}
	defer rows.Close()

	var partUUIDs []uuid.UUID
	for rows.Next() {
		var pu uuid.UUID
		if err := rows.Scan(&pu); err != nil {
			return model.OrderDto{}, fmt.Errorf("scan part_uuid: %w", err)
		}
		partUUIDs = append(partUUIDs, pu)
	}
	if err := rows.Err(); err != nil {
		return model.OrderDto{}, fmt.Errorf("iterate order_parts: %w", err)
	}

	order := repoConverter.RepoToServiceModel(repoOrder)
	order.PartUuids = partUUIDs
	return order, nil
}
