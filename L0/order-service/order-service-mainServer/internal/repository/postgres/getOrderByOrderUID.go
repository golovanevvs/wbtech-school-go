package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/model"
)

// ErrNoExists is returned when a order uid is not found in database.
var ErrOrderNotFound = errors.New("order not found in database")

func (p *Postgres) GetOrderByOrderUID(ctx context.Context, orderUID string) (*model.Order, error) {
	log := p.logger.With().
		Str("order_uid", orderUID).
		Str("method", "GetOrderByOrderUID").
		Logger()

	log.Debug().Msg("starting GetOrderByOrderUID")

	type orderRow struct {
		model.Order
		Delivery model.Delivery `db:"delivery"`
		Payment  model.Payment  `db:"payment"`
	}

	var row orderRow
	err := p.db.GetContext(ctx, &row, `

		SELECT 
			o.*,
			d.name as "delivery.name",
			d.phone as "delivery.phone",
			d.zip as "delivery.zip",
			d.city as "delivery.city",
			d.address as "delivery.address",
			d.region as "delivery.region",
			d.email as "delivery.email",
			p.transaction as "payment.transaction",
			p.request_id as "payment.request_id",
			p.currency as "payment.currency",
			p.provider as "payment.provider",
			p.amount as "payment.amount",
			p.payment_dt as "payment.payment_dt",
			p.bank as "payment.bank",
			p.delivery_cost as "payment.delivery_cost",
			p.goods_total as "payment.goods_total",
			p.custom_fee as "payment.custom_fee"
		FROM orders o
		LEFT JOIN delivery d ON d.order_id = o.id
		LEFT JOIN payment p ON p.order_id = o.id
		WHERE o.order_uid = $1

	`, orderUID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Warn().Msg("order not found in database")
			return nil, ErrOrderNotFound
		}
		log.Error().Err(err).Msg("failed to get order")
		return nil, err
	}

	order := row.Order
	order.Delivery = row.Delivery
	order.Payment = row.Payment

	err = p.db.SelectContext(ctx, &order.Items, `

		SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		FROM items
		WHERE order_id = $1

	`, order.ID)
	if err != nil {
		log.Error().Err(err).
			Int("order_id", order.ID).Msg("failed to get items for order")
		return nil, err
	}

	log.Debug().
		Int("items_count", len(order.Items)).
		Msg("order successfully retrieved")
	return &order, nil
}
