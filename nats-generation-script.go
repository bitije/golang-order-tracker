package main

import (
	"log"

	"github.com/nats-io/stan.go"
)

func main() {
	sc, err := stan.Connect("test-cluster", "Publisher", stan.NatsURL("nats://localhost:4223"))
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()
	messages := make([]string, 3)
	messages[0] = `
    {
        "order_id": "3c51e153-0c79-4613-9ee0-3302bb6d8343",
        "track_number": "vmurweoirw",
        "entry": "fdjfks",
        "delivery": {
            "name": "Test Testov",
            "phone": "+9720000000",
            "zip": "2639809",
            "city": "Kiryat Mozkin",
            "address": "Ploshad Mira 15",
            "region": "Kraiot",
            "email": "test@gmail.com"
        },
        "payment": {
            "transaction": "b563feb7b2b84b6test",
            "request_id": "",
            "currency": "USD",
            "provider": "wbpay",
            "amount": 1817,
            "payment_dt": 1637907727,
            "bank": "alpha",
            "delivery_cost": 1500,
            "goods_total": 317,
            "custom_fee": 0
        },
        "items": [
        {
            "chrt_id": 9934930,
            "track_number": "fjksdlfs",
            "price": 453,
            "rid": "ab4219087a764ae0btest",
            "name": "Mascaras",
            "sale": 30,
            "size": "0",
            "total_price": 317,
            "nm_id": 2389212,
            "brand": "Vivienne Sabo",
            "status": 202
        },
        {
            "chrt_id": 9934930,
            "track_number": "fjksdlfs",
            "price": 453,
            "rid": "ab4219087a764ae0btest",
            "name": "Mascaras",
            "sale": 30,
            "size": "0",
            "total_price": 317,
            "nm_id": 2389212,
            "brand": "Vivienne Sabo",
            "status": 202
        }
        ],
        "locale": "en",
        "internal_signature": "",
        "customer_id": "test",
        "delivery_service": "meest",
        "shardkey": "9",
        "sm_id": 99,
        "date_created": "2021-11-26T06:22:19Z",
        "oof_shard": "1"
    }
    `
	messages[1] = `
    {
        "order_id": "1a9c4c15-2ffe-4633-9c51-2a51ffc43bac",
        "track_number": "WBILMTRACK",
        "entry": "WBIL",
        "delivery": {
            "name": "Test Testov",
            "phone": "+9720000000",
            "zip": "2639809",
            "city": "Kiryat Mozkin",
            "address": "Ploshad Mira 15",
            "region": "Kraiot",
            "email": "test@gmail.com"
        },
        "payment": {
            "transaction": "b563feb7b2b84b6",
            "request_id": "",
            "currency": "USD",
            "provider": "wbpay",
            "amount": 1816,
            "payment_dt": 163790727,
            "bank": "alpha",
            "delivery_cost": 500,
            "goods_total": 417,
            "custom_fee": 0
        },
        "items": [
        {
            "chrt_id": 8934530,
            "track_number": "jfrueow",
            "price": 455,
            "rid": "ab4219087a764ae0b",
            "name": "Jopa",
            "sale": 30,
            "size": "0",
            "total_price": 317,
            "nm_id": 2389212,
            "brand": "DZHOPA",
            "status": 202
        }
        ],
        "locale": "en",
        "internal_signature": "",
        "customer_id": "test",
        "delivery_service": "meest",
        "shardkey": "9",
        "sm_id": 99,
        "date_created": "2021-11-26T06:22:19Z",
        "oof_shard": "1"
    }
    `
	messages[2] = `
    {
        "order_id": "dd0fd574-e339-42b4-93d2-471744774d3c",
        "track_number": "dfdsjksklfd",
        "entry": "fjrwiiei",
        "delivery": {
            "name": "Test",
            "phone": "+9720000000",
            "zip": "2639809",
            "city": "Kiryat Mozkin",
            "address": "Ploshad Mira 15",
            "region": "Kraiot",
            "email": "test@gmail.com"
        },
        "payment": {
            "transaction": "b563b2bb6",
            "request_id": "",
            "currency": "USD",
            "provider": "pay",
            "amount": 18227,
            "payment_dt": 1907727,
            "bank": "alpha",
            "delivery_cost": 100,
            "goods_total": 37,
            "custom_fee": 0
        },
        "items": [
        {
            "chrt_id": 49,
            "track_number": "vmnurw",
            "price": 3,
            "rid": "ab4219087a764ae0btest",
            "name": "Mascaras",
            "sale": 30,
            "size": "0",
            "total_price": 3,
            "nm_id": 2389212,
            "brand": "Sabo",
            "status": 2
        }
        ],
        "locale": "en",
        "internal_signature": "",
        "customer_id": "tt",
        "delivery_service": "st",
        "shardkey": "5",
        "sm_id": 9,
        "date_created": "2021-11-26T06:22:19Z",
        "oof_shard": "1"
    }
    `
	for idx, msg := range messages {
		log.Printf("Publishing a message: %d", idx)
		err = sc.Publish("foo", []byte(msg))
		if err != nil {
			log.Fatal(err)
		}
	}
}
