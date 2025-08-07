package postgres

import (
	"context"

	"github.com/golovanevvs/wbtech-school-go/L0/order-service/order-service-mainServer/internal/model"
)

func (p *Postgres) AddOrder(ctx context.Context, order model.Order) error {
	log := p.logger.With().
		Str("order_uid", order.OrderUID).
		Str("method", "AddOrder").
		Logger()

	log.Debug().Msg("starting AddOrder")

	log.Trace().Msg("starting transaction")
	tx, err := p.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to begin transaction")
		return err
	}
	defer func() {
		if err != nil {
			log.Error().Err(err).Msg("rolling back transaction")
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Error().Err(rbErr).Msg("failed to rollback transaction")
			}
		}
	}()

	log.Trace().Msg("inserting order record")
	row := tx.QueryRowContext(ctx, `
	INSERT INTO orders
		(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_chard)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id;
	`, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)

	var orderID int
	if err := row.Scan(&orderID); err != nil {
		log.Error().Err(err).Msg("failed to insert order")
		return err
	}

	log.Trace().Msg("inserting delivery record")
	_, err = tx.ExecContext(ctx, `
	INSERT INTO delivery
		(name, phone, zip, city, address, region, email, order_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8);
	`, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email, orderID)
	if err != nil {
		log.Error().Err(err).
			Int("order_id", orderID).
			Msg("failed to insert delivery")
		return err
	}

	log.Trace().Msg("inserting payment record")
	_, err = tx.ExecContext(ctx, `
	INSERT INTO payment
		(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee, order_id)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee, orderID)
	if err != nil {
		log.Error().Err(err).
			Int("order_id", orderID).
			Msg("failed to insert payment")
		return err
	}

	if len(order.Items) > 0 {
		log.Trace().
			Int("items_count", len(order.Items)).
			Msg("inserting items")
		for i := range order.Items {
			order.Items[i].OrderID = orderID
		}
		_, err = tx.NamedExecContext(ctx, `
		INSERT INTO items
			(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_id)
		VALUES
			(:chrt_id, :track_number, :price, :rid, :name, :sale, :size, :total_price, :nm_id, :brand, :status, :order_id);
		`, order.Items)
		if err != nil {
			log.Error().Err(err).
				Int("order_id", orderID).
				Msg("failed to insert items")
			return err
		}
	}

	log.Trace().Msg("commiting transaction")
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).
			Msg("failed to commit transaction")
		return err
	}

	log.Debug().
		Int("order_id", orderID).
		Msg("order successfully added to database")

	return nil
}
