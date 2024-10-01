package database

import (
	"L0/config"
	"L0/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

func Connect(cfg *config.DB) (*pgxpool.Pool, error) {
	conn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	return pgxpool.Connect(context.Background(), conn)
}
func NewDB(pool *pgxpool.Pool) *Postgres {
	return &Postgres{
		pool: pool,
	}
}
func (db *Postgres) CreateTable() error {
	_, err := db.pool.Exec(context.Background(), `
	CREATE SCHEMA IF NOT EXISTS orders;
	CREATE TABLE IF NOT EXISTS orders.main (
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
func (db *Postgres) SaveOrder(order models.OrderJSON) error {
	_, err := db.pool.Exec(context.Background(), `
	INSERT INTO orders.main (order_uid, track_number, entry, delivery_info, payment_info, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`, order.OrderUid, order.TrackNumber, order.Entry, order.Delivery, order.Payments, order.Items, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OOFShard)
	return err
}
func (db *Postgres) GetAllOrders() (orders []models.OrderJSON, err error) {
	res, err := db.pool.Query(context.Background(), `SELECT * FROM orders.main`)
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
			&executeOrder.SmId,
			&executeOrder.DateCreated,
			&executeOrder.OOFShard,
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
func (db *Postgres) GetOrderByUID(uid string) (models.OrderJSON, error) {
	var executeOrder models.OrderJSON
	err := db.pool.QueryRow(context.Background(), `SELECT * FROM orders.main WHERE order_uid=$1`, uid).Scan(
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
		&executeOrder.SmId,
		&executeOrder.DateCreated,
		&executeOrder.OOFShard,
	)
	if err != nil {
		return executeOrder, err
	}
	return executeOrder, nil
}
