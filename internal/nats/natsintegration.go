package nats

import (
	"log"

	"github.com/nats-io/stan.go"
	"github.com/vojdelenie/task0/internal/db"
)

type NatsConfig struct {
	Connection  stan.Conn
	ClientID    string
	ClusterID   string
	NatsURL     string
	DurableName string
	Subject     string
}

var Instance NatsConfig

func Run() {
	Instance.ClientID = "Consumer"
	Instance.ClusterID = "test-cluster"
	Instance.NatsURL = "nats://localhost:4223"
	Instance.DurableName = "test-durable-name"
	Instance.Subject = "foo"
	err := Instance.connect()
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	} else {
		log.Printf("Connection to NATS established")
	}
	Instance.subscribeAndProcessMessages()
}

func (n *NatsConfig) connect() error {
	conn, err := stan.Connect(n.ClusterID, n.ClientID, stan.NatsURL(n.NatsURL))
	if err != nil {
		return err
	}
	n.Connection = conn
	return err
}

func (n *NatsConfig) subscribeAndProcessMessages() {
	_, err := n.Connection.Subscribe(n.Subject, func(msg *stan.Msg) {
		insertResult := make(chan error)
		orderId := make(chan string)
		go db.Instance.InsertOrder(string(msg.Data), insertResult, orderId)
		err := <-insertResult
		if err != nil {
			log.Printf("Error inserting order: %v", err)
		} else {
			go db.Cache.Set(<-orderId, string(msg.Data))
		}
	}, stan.DurableName(n.DurableName))
	if err != nil {
		log.Printf("Error subscribing to subject %s: %v", n.Subject, err)
	}
}
