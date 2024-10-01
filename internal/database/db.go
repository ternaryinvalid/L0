package database

import (
	"L0/config"
	"L0/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type OrdersRepository struct {
	pool   *pgxpool.Pool
	schema string
}

func Connect(cfg *config.DB, ctx context.Context) (*pgxpool.Pool, error) {

	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	return pgxpool.Connect(ctx, conn)
}

func NewDB(pool *pgxpool.Pool, schema string) *OrdersRepository {
	return &OrdersRepository{
		pool:   pool,
		schema: schema,
	}
}

func (db *OrdersRepository) CreateSchemaAndTable(ctx context.Context) error {

	_, err := db.pool.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS orders.orders (
		order_uid VARCHAR(255) PRIMARY KEY,
        track_number VARCHAR(255),
        entry VARCHAR(255),
        delivery_info JSONB,
        payment_info JSONB,
        items JSONB,
        locale VARCHAR(255),
        internal_signature VARCHAR(255),
        customer_id VARCHAR(255),
        delivery_service VARCHAR(255),
        shardkey VARCHAR(255),
        sm_id INTEGER,
        date_created VARCHAR(255),
        oof_shard VARCHAR(255)
	)
	`)
	return err
}

func (db *OrdersRepository) SaveOrder(order models.OrderJSON, ctx context.Context) error {
	_, err := db.pool.Exec(ctx, `
	INSERT INTO orders.orders VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, order.OrderUid, order.TrackNumber, order.Entry, order.Delivery, order.Payments, order.Items, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.Sm_id, order.Date_created, order.OOF_shard)
	return err
}

func (db *OrdersRepository) GetAllOrders(ctx context.Context) (orders []models.OrderJSON, err error) {
	res, err := db.pool.Query(ctx, `SELECT * FROM orders.orders`)

	if err != nil {
		fmt.Printf("Query error: %v", err)
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var executeOrder models.OrderJSON
		err := res.Scan(
			&executeOrder.OrderUid,
			&executeOrder.TrackNumber,
			&executeOrder.Entry,
			&executeOrder.Delivery,
			&executeOrder.Payments,
			&executeOrder.Items,
			&executeOrder.Locale,
			&executeOrder.InternalSignature,
			&executeOrder.CustomerId,
			&executeOrder.DeliveryService,
			&executeOrder.Shardkey,
			&executeOrder.Sm_id,
			&executeOrder.Date_created,
			&executeOrder.OOF_shard,
		)

		if err != nil {
			fmt.Printf("Error executing order: %v\n", err)
		}

		orders = append(orders, executeOrder)
	}

	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("error executing order: %v", err)
	}

	fmt.Println("Orders downloaded to cache")
	return orders, nil
}

func (db *OrdersRepository) GetOrderByUID(uid string, ctx context.Context) (models.OrderJSON, error) {

	var executeOrder models.OrderJSON
	err := db.pool.QueryRow(ctx, `SELECT * FROM orders.orders WHERE order_uid=$1`, uid).Scan(
		&executeOrder.OrderUid,
		&executeOrder.TrackNumber,
		&executeOrder.Entry,
		&executeOrder.Delivery,
		&executeOrder.Payments,
		&executeOrder.Items,
		&executeOrder.Locale,
		&executeOrder.InternalSignature,
		&executeOrder.CustomerId,
		&executeOrder.DeliveryService,
		&executeOrder.Shardkey,
		&executeOrder.Sm_id,
		&executeOrder.Date_created,
		&executeOrder.OOF_shard,
	)

	if err != nil {
		return executeOrder, err
	}

	return executeOrder, nil
}
