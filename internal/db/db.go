package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type DB struct {
	*sql.DB
}

var Instance DB

const (
	dbDriver   = "postgres"
	dbUser     = "postgres"
	dbPassword = "example"
	dbName     = "gopg"
)

func Run() error {
	err := Instance.Connect()
	if err != nil {
		log.Fatalf("Error connecting to DB: %v\n", err)
		return err
	} else {
		log.Printf("Connection to DB established")
	}
	err = Instance.CreateTables()
	if err != nil {
		log.Fatalf("Error creating tables: %v\n", err)
		return err
	} else {
		log.Printf("Tables created")
	}
	return nil
}

func (db *DB) Connect() error {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)
	database, err := sql.Open(dbDriver, connStr)
	if err != nil {
		return err
	}
	if err := database.Ping(); err != nil {
		return err
	}
	db.DB = database
	return nil
}

func (db *DB) InsertOrder(jsonData string, orderInserted chan error, res chan string) {
	var order Order
	err := json.Unmarshal([]byte(jsonData), &order)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v\n", err)
		orderInserted <- err
		return
	}

	order.Delivery.DeliveryID = uuid.New()
	_, err = db.Exec(`
        INSERT INTO delivery (delivery_id, name, phone, zip, city, address, region, email)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.Delivery.DeliveryID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		log.Printf("Error inserting delivery: %v\n", err)
		orderInserted <- err
		return
	}

	order.Payment.PaymentID = uuid.New()
	_, err = db.Exec(`
        INSERT INTO payment (payment_id, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.Payment.PaymentID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, float64(order.Payment.Amount)/100,
		order.Payment.PaymentDt, order.Payment.Bank, float64(order.Payment.DeliveryCost)/100,
		float64(order.Payment.GoodsTotal)/100, float64(order.Payment.CustomFee)/100,
	)
	if err != nil {
		log.Printf("Error inserting payment: %v\n", err)
		orderInserted <- err
		return
	}

	order.OrderID = uuid.New()
	_, err = db.Exec(`
        INSERT INTO orders (order_id, track_number, entry, delivery, payment,
            locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`,
		order.OrderID, order.TrackNumber, order.Entry, order.Delivery.DeliveryID,
		order.Payment.PaymentID, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.ShardKey,
		order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		log.Printf("Error inserting order: %v\n", err)
		orderInserted <- err
		return
	}

	for i := 0; i < (len(order.Items)); i++ {
		order.Items[i].ItemID = uuid.New()
		_, err = db.Exec(`
            INSERT INTO items (item_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			order.Items[i].ItemID, order.Items[i].ChrtID, order.Items[i].TrackNumber,
			order.Items[i].Price, order.Items[i].RID, order.Items[i].Name,
			order.Items[i].Sale, order.Items[i].Size, order.Items[i].TotalPrice,
			order.Items[i].NmID, order.Items[i].Brand, order.Items[i].Status,
		)
		if err != nil {
			log.Printf("Error inserting item: %v\n", err)
			orderInserted <- err
			return
		}

		_, err = db.Exec(`
        INSERT INTO order_items (order_id, item_id)
        VALUES ($1, $2)`,
			order.OrderID, order.Items[i].ItemID,
		)
		if err != nil {
			log.Printf("Error inserting relationship: %v\n", err)
			orderInserted <- err
			return
		}
	}

	log.Printf("Inserted order: %s\n", order.OrderID)
	orderInserted <- nil
	res <- order.OrderID.String()
}

func (db *DB) GetOrder(id uuid.UUID) (*Order, error) {
	query := `
        SELECT
        orders.order_id,
        orders.track_number,
        orders.entry,

        delivery.delivery_id,
        delivery.name,
        delivery.phone,
        delivery.zip,
        delivery.city,
        delivery.address,
        delivery.region,
        delivery.email,

        payment.payment_id,
        payment.transaction,
        payment.request_id,
        payment.currency,
        payment.provider,
        payment.amount,
        payment.payment_dt,
        payment.bank,
        payment.delivery_cost,
        payment.goods_total,
        payment.custom_fee,

        items.item_id,
        items.chrt_id,
        items.track_number AS item_track_number,
        items.price,
        items.rid,
        items.name AS item_name,
        items.sale,
        items.size,
        items.total_price,
        items.nm_id,
        items.brand,
        items.status,

        orders.locale,
        orders.internal_signature,
        orders.customer_id,
        orders.delivery_service,
        orders.shardkey,
        orders.sm_id,
        orders.date_created,
        orders.oof_shard
        FROM orders 
        LEFT JOIN delivery ON orders.delivery = delivery.delivery_id
        LEFT JOIN payment ON orders.payment = payment.payment_id
        LEFT JOIN order_items ON orders.order_id = order_items.order_id
        LEFT JOIN items ON order_items.item_id = items.item_id
        WHERE orders.order_id = $1
    `
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var order Order
	var delivery Delivery
	var payment Payment
	var items []Item

	for rows.Next() {
		var item Item
		err := rows.Scan(
			&order.OrderID, &order.TrackNumber, &order.Entry,

			&delivery.DeliveryID, &delivery.Name, &delivery.Phone, &delivery.Zip,
			&delivery.City, &delivery.Address, &delivery.Region, &delivery.Email,

			&payment.PaymentID, &payment.Transaction, &payment.RequestID, &payment.Currency,
			&payment.Provider, &payment.Amount, &payment.PaymentDt, &payment.Bank,
			&payment.DeliveryCost, &payment.GoodsTotal, &payment.CustomFee,

			&item.ItemID, &item.ChrtID, &item.TrackNumber, &item.Price,
			&item.RID, &item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status,

			&order.Locale, &order.InternalSignature, &order.CustomerID,
			&order.DeliveryService, &order.ShardKey,
			&order.SmID, &order.DateCreated, &order.OofShard,
		)
		if err != nil {
			return nil, err
		}
		order.Delivery = delivery
		order.Payment = payment
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	order.Items = items
	log.Printf("Got order %v from DB.\n", order.OrderID)
	return &order, nil
}

func (db *DB) CreateTables() error {
	query := `
    CREATE TABLE IF NOT EXISTS delivery (
        delivery_id UUID PRIMARY KEY,
        name VARCHAR(50),
        phone VARCHAR(50),
        zip VARCHAR(50),
        city VARCHAR(50),
        address VARCHAR(255),
        region VARCHAR(50),
        email VARCHAR(50)
    );

    CREATE TABLE IF NOT EXISTS payment (
        payment_id UUID PRIMARY KEY,
        transaction VARCHAR(50),
        request_id VARCHAR(50),
        currency VARCHAR(50),
        provider VARCHAR(50),
        amount DECIMAL(10, 2),
        payment_dt INT,
        bank VARCHAR(50),
        delivery_cost DECIMAL(10, 2),
        goods_total DECIMAL(10, 2),
        custom_fee DECIMAL(10, 2)
    );

    CREATE TABLE IF NOT EXISTS items (
        item_id UUID PRIMARY KEY,
        chrt_id INT,
        track_number VARCHAR(50),
        price INT,
        rid VARCHAR(50),
        name VARCHAR(255),
        sale INT,
        size VARCHAR(50),
        total_price INT,
        nm_id INT,
        brand VARCHAR(50),
        status INT
    );

    CREATE TABLE IF NOT EXISTS orders (
        order_id UUID PRIMARY KEY,
        track_number VARCHAR(50),
        entry VARCHAR(50),
        delivery UUID REFERENCES delivery(delivery_id),
        payment UUID REFERENCES payment(payment_id),
        locale VARCHAR(10),
        internal_signature VARCHAR(50),
        customer_id VARCHAR(50),
        delivery_service VARCHAR(50),
        shardkey VARCHAR(50),
        sm_id INT,
        date_created TIMESTAMP,
        oof_shard VARCHAR(50)
    );
    CREATE TABLE IF NOT EXISTS order_items (
        order_id UUID REFERENCES orders(order_id),
        item_id UUID REFERENCES items(item_id),
        PRIMARY KEY (order_id, item_id)
    );
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Print(err)
	}
	return err
}
