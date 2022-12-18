package rabbitmq

import (
	"fmt"
	"log"
	"time"

	"github.com/micro-kit/micro-common/config"
	"github.com/micro-kit/micro-common/etcdcli"
	"github.com/micro-kit/micro-common/logger"
	"github.com/streadway/amqp"
)

var (
	rabbitmqHandle *RabbitmqHandler
	rabbitmqConfig *config.RabbitmqConfig
)

/* rabbitmq 操作 */
type RabbitmqHandler struct {
	client        *amqp.Connection
	cfg           *config.RabbitmqConfig
	notifyReConns []NotifyReConn
}

func GetRabbitmqHandle() *RabbitmqHandler {
	if rabbitmqHandle == nil {
		initRabbitmq()
	}
	return rabbitmqHandle
}

func init() {
	initRabbitmq()
}

func initRabbitmq() {
	log.Println("开始初始化rabbitmq连接")
	var err error
	err = config.GetRabbitmqConfig(etcdcli.EtcdCli, func(cfg *config.RabbitmqConfig) {
		rabbitmqConfig = cfg
		rabbitmqHandle, err = NewRabbitmqHandler(cfg)
		if err != nil {
			logger.Logger.Panicw("Creating rabbitmq connection errors", "err", err)
		}
	})
	if err != nil {
		logger.Logger.Panicw("Get rabbitmq configuration error", "err", err)
	}
}

func NewRabbitmqHandler(cfg *config.RabbitmqConfig) (*RabbitmqHandler, error) {
	mqClient, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Vhost))
	if err != nil {
		return nil, err
	}

	rh := &RabbitmqHandler{
		client:        mqClient,
		cfg:           cfg,
		notifyReConns: make([]NotifyReConn, 0),
	}
	// 处理连接错误
	go rh.notifyClose(mqClient.NotifyClose(make(chan *amqp.Error)))

	return rh, nil
}

// 处理连接错误
func (rh *RabbitmqHandler) notifyClose(chanClose chan *amqp.Error) {
	err := <-chanClose
	if err == nil {
		return
	}
	log.Println("rabbitmq连接错误", "err", err)
	log.Println("rabbitmq开始重新连接")
	for {
		mqClient, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rh.cfg.User, rh.cfg.Password, rh.cfg.Host, rh.cfg.Port, rh.cfg.Vhost))
		if err != nil {
			log.Println("rabbitmq重连失败", err)
			time.Sleep(3 * time.Second)
		} else {
			rh.client = mqClient
			break
		}
	}

	// 通知重新订阅
	notifyReConns := rh.notifyReConns
	rh.notifyReConns = make([]NotifyReConn, 0) // 防止订阅放重新定义被清空
	for _, f := range notifyReConns {
		go f()
	}

	// 处理连接错误
	go rh.notifyClose(rh.client.NotifyClose(make(chan *amqp.Error)))
}

// 发送mq消息
func (rh *RabbitmqHandler) SendMessage(exchange, exchangeType, routingKey, body string, bodyType string) error {
	ch, err := rh.client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclarePassive(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     bodyType,
			ContentEncoding: "utf-8",
			Body:            []byte(body),
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

type NotifyReConn func() // 连接通知
type MessageCallBack func([]byte, string) bool

// 订阅消息
func (rh *RabbitmqHandler) ReceiveMessage(exchange, exchangeType, queueName, key string, callBack MessageCallBack, notifyReConn NotifyReConn) error {
	ch, err := rh.client.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// 添加通知回调
	if notifyReConn != nil {
		rh.notifyReConns = append(rh.notifyReConns, notifyReConn)
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclarePassive(
		exchange,
		exchangeType,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err = ch.QueueBind(
		queue.Name,
		key,
		exchange,
		false,
		nil,
	); err != nil {
		return err
	}

	deliveries, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for d := range deliveries {
		ack := callBack(d.Body, d.RoutingKey)
		d.Ack(ack)
	}

	return nil
}

// 销毁
func (rh *RabbitmqHandler) Destroy() {
	if rh.client == nil {
		return
	}
	rh.client.Close()
}
