package rabbitmq

import (
	"fmt"

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
	client *amqp.Connection
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

	return &RabbitmqHandler{
		client: mqClient,
	}, nil
}
