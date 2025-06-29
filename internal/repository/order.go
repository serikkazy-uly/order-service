package repository

import (
	"database/sql"
	"order-service/internal/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	apperrors "order-service/internal/errors" // кастомнаые ошибки
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(order *models.Order) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Вставка основного заказа
	_, err = tx.NamedExec(`
        INSERT INTO orders (order_uid, track_number, entry, locale, 
                          internal_signature, customer_id, delivery_service, 
                          shardkey, sm_id, date_created, oof_shard)
        VALUES (:order_uid, :track_number, :entry, :locale, 
                :internal_signature, :customer_id, :delivery_service, 
                :shardkey, :sm_id, :date_created, :oof_shard)
    `, order)

	if err != nil {
		return err
	}

	// Вставка доставки
	deliveryQuery := `
        INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `
	_, err = tx.Exec(deliveryQuery, order.OrderUID, order.Delivery.Name,
		order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return err
	}

	// Вставка платежа
	paymentQuery := `
        INSERT INTO payments (order_uid, transaction, request_id, currency, provider,
                            amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
    `
	_, err = tx.Exec(paymentQuery, order.OrderUID, order.Payment.Transaction,
		order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider,
		order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return err
	}

	// Вставка товаров
	for _, item := range order.Items {
		itemQuery := `
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name,
                             sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `
		_, err = tx.Exec(itemQuery, order.OrderUID, item.ChrtID, item.TrackNumber,
			item.Price, item.Rid, item.Name, item.Sale, item.Size,
			item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *OrderRepository) GetOrder(orderUID string) (*models.Order, error) {
	var order models.Order

	// Получаем основную информацию о заказе
	err := r.db.Get(&order, `
        SELECT order_uid, track_number, entry, locale, internal_signature,
               customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders WHERE order_uid = $1
    `, orderUID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrOrderNotFound
		}

		return nil, err
	}

	// Получаем информацию о доставке
	err = r.db.Get(&order.Delivery, `
        SELECT name, phone, zip, city, address, region, email
        FROM deliveries WHERE order_uid = $1
    `, orderUID)

	if err != nil {
		return nil, err
	}

	// Получаем информацию о платеже
	err = r.db.Get(&order.Payment, `
        SELECT transaction, request_id, currency, provider, amount,
               payment_dt, bank, delivery_cost, goods_total, custom_fee
        FROM payments WHERE order_uid = $1
    `, orderUID)

	if err != nil {
		return nil, err
	}

	// Получаем товары
	err = r.db.Select(&order.Items, `
        SELECT chrt_id, track_number, price, rid, name, sale, size,
               total_price, nm_id, brand, status
        FROM items WHERE order_uid = $1
    `, orderUID)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order

	rows, err := r.db.Query(`
        SELECT DISTINCT order_uid FROM orders
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderUIDs []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		orderUIDs = append(orderUIDs, uid)
	}

	for _, uid := range orderUIDs {
		order, err := r.GetOrder(uid)
		if err != nil {
			continue // Пропускаем поврежденные записи
		}
		orders = append(orders, *order)
	}

	return orders, nil
}
