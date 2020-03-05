package nsqlib

import (
	"os"

	"github.com/nsqio/go-nsq"
)

/* nsq 访问客户端 */

// NewProducer 获取一个nsq生产者客户端 - 注意使用完 producer.Stop()
func NewProducer() (*nsq.Producer, error) {
	url := os.Getenv("NSQ_URL")
	producer, err := nsq.NewProducer(url, nsq.NewConfig())
	if err != nil {
		return nil, err
	}
	mode := os.Getenv("DEFAULT_MODE")
	if mode == "dev" {
		producer.SetLoggerLevel(nsq.LogLevelDebug)
	} else if mode == "test" {
		producer.SetLoggerLevel(nsq.LogLevelInfo)
	} else {
		producer.SetLoggerLevel(nsq.LogLevelWarning)
	}
	err = producer.Ping()
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// NewConsumer 创建一个消费者
func NewConsumer(topic string, channel string, handle nsq.Handler) error {
	config := nsq.NewConfig()
	config.MaxInFlight = 9
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return err
	}
	mode := os.Getenv("DEFAULT_MODE")
	if mode == "dev" {
		consumer.SetLoggerLevel(nsq.LogLevelDebug)
	} else if mode == "test" {
		consumer.SetLoggerLevel(nsq.LogLevelInfo)
	} else {
		consumer.SetLoggerLevel(nsq.LogLevelWarning)
	}
	consumer.AddHandler(handle)
	url := os.Getenv("NSQ_URL")
	err = consumer.ConnectToNSQD(url)
	if err != nil {
		return err
	}
	return nil
}
