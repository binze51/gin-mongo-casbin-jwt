package pulsarmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

type ComsumerConf struct {
	URL              string        `yaml:"PULSAR_URL"     envDefault:"pulsar://localhost:6650"`
	OpTimeOut        time.Duration `yaml:"OPTIMEOUT"     envDefault:"30 * time.Second"`
	ConnTimeOut      time.Duration `yaml:"CONNTIMEOUT"     envDefault:"30 * time.Second"`
	AuthToken        string
	CallServerName   string
	SubscriptionName string
	TopicName        string
	Type             int //消费模式
}

type Consumer struct {
	Consumer pulsar.Consumer
}

// CreateComsumer - 消费者
func CreateComsumer(conf *ComsumerConf) (*Consumer, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               conf.URL,
		OperationTimeout:  conf.OpTimeOut,
		ConnectionTimeout: conf.ConnTimeOut,
	})

	if err != nil {
		log.Printf("Could not instantiate Pulsar client: %v", err)
		return nil, err
	}

	// create a new producer instance
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            conf.TopicName,
		SubscriptionName: conf.SubscriptionName,
		Type:             pulsar.SubscriptionType(conf.Type),
	})
	if err != nil {
		log.Fatal(err)
	}
	return &Consumer{Consumer: consumer}, nil

}

func (c *Consumer) ListenerMessge() {
	log.Println("Listeneing for events")
	for {
		msg, err := c.Consumer.Receive(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		payload := msg.Payload()
		var message interface{}
		err = json.Unmarshal(payload, &message)
		if err != nil {
			fmt.Errorf("Error saving target: %v", err)
			return
		}
		println(message)
		// err = e.Repo.SaveTarget(targets)
		// if err != nil {
		// 	fmt.Errorf("Error saving target: %v", err)
		// 	return
		// }
		// 手动确认接收的消息
		c.Consumer.AckID(msg.ID())
		fmt.Println("消息已确认接收.")
	}
}
