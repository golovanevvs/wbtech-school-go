package postgres

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service_main-server/internal/model"
)

func (p *Postgres) GetOrdersForCache(ctx context.Context) ([]model.Order, error) {
	log := p.logger.With().
		Str("method", "GetOrdersForCache").
		Logger()

	log.Debug().Msg("starting GetOrdersForCache")

	type orderRow struct {
		model.Order
		Delivery model.Delivery `db:"delivery"`
		Payment  model.Payment  `db:"payment"`
	}

	var orders []model.Order
	var rows []orderRow
	err := p.db.SelectContext(ctx, &rows, `

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
		ORDER BY o.id DESC 
		LIMIT 10
	
	`)
	if err != nil {
		log.Error().Err(err).Msg("failed to get orders")
		return nil, err
	}

	for _, row := range rows {
		order := row.Order
		order.Delivery = row.Delivery
		order.Payment = row.Payment

		err := p.db.SelectContext(ctx, &order.Items, `
			
			SELECT
				chrt_id, track_number, price, rid, name, sale,
				size, total_price, nm_id, brand, status
			FROM items
			WHERE order_id = $1

		`, order.ID)
		if err != nil {
			log.Error().Err(err).Msg("failed to get items for order")
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
