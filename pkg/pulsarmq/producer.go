package pulsarmq

import (
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

type ProducerConf struct {
	URL              string        `yaml:"PULSAR_URL"     envDefault:"pulsar://localhost:6650"`
	OpTimeOut        time.Duration `yaml:"OPTIMEOUT"     envDefault:"30 * time.Second"`
	ConnTimeOut      time.Duration `yaml:"CONNTIMEOUT"     envDefault:"30 * time.Second"`
	AuthToken        string
	SubscriptionName string
	TopicName        string
}

type Producer struct {
	Producer pulsar.Producer
}

// CreateProducer - 生产者
func CreateProducer(conf *ProducerConf) (pulsar.Producer, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               conf.URL,
		OperationTimeout:  conf.OpTimeOut,
		ConnectionTimeout: conf.ConnTimeOut,
	})

	if err != nil {
		log.Printf("Could not instantiate Pulsar client: %v", err)
		return nil, err
	}

	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: conf.TopicName,
		Name:  conf.SubscriptionName,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Producer is ready")
	return producer, nil
}
