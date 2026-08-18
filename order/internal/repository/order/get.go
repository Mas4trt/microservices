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
	var (
		repoOrder repoModel.OrderDto
		partUUIDs []uuid.UUID
	)

	err := r.pool.QueryRow(ctx, `
		SELECT
			o.order_uuid,
			o.user_uuid,
			o.total_price,
			o.transaction_uuid,
			o.payment_method,
			o.status,
			COALESCE(
				array_agg(op.part_uuid) FILTER (WHERE op.part_uuid IS NOT NULL),
				'{}'
			)
		FROM orders o
		LEFT JOIN order_parts op ON op.order_uuid = o.order_uuid
		WHERE o.order_uuid = $1
		GROUP BY
			o.order_uuid,
			o.user_uuid,
			o.total_price,
			o.transaction_uuid,
			o.payment_method,
			o.status
	`, orderUUID).Scan(
		&repoOrder.OrderUUID,
		&repoOrder.UserUUID,
		&repoOrder.TotalPrice,
		&repoOrder.TransactionUUID,
		&repoOrder.PaymentMethod,
		&repoOrder.Status,
		&partUUIDs,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.OrderDto{}, model.ErrOrderNotFound
		}
		return model.OrderDto{}, fmt.Errorf("query order: %w", err)
	}

	order := repoConverter.RepoToServiceModel(repoOrder)
	order.PartUuids = partUUIDs

	return order, nil
}
