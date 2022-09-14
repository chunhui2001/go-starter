package grabbit

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

var (
	logger  *logrus.Entry
	server  string
	channel *amqp.Channel
)

type GRabbitConf struct {
	Enable      bool     `mapstructure:"RABBITMQ_ENABLE"`
	AutoConsume bool     `mapstructure:"RABBITMQ_AUTO_CONSUME"`
	Server      string   `mapstructure:"RABBITMQ_SERVER"`
	Queues      []string `mapstructure:"RABBITMQ_QUEUES"` // 多个逗号分隔
	GroupId     string   `mapstructure:"RABBITMQ_GROUP_ID"`
}

func Init(conf *GRabbitConf, log *logrus.Entry) {

	logger = log
	server = conf.Server

	conn, err := amqp.Dial(server)

	if err != nil {
		logger.Errorf("RabbitMQ-Failed-to-Connect: Server=%s, errorMessage=%s", server, err.Error())
		return
	}

	ch, err := conn.Channel()

	if err != nil {
		logger.Errorf("RabbitMQ-Failed-to-Open-a-Channel: Server=%s, errorMessage=%s", server, err.Error())
		return
	}

	channel = ch

	logger.Infof("RabbitMQ-Connect-Success: Server=%s", server)

	for _, q := range conf.Queues {
		if err := QueueDeclare(q); err != nil {
			logger.Errorf("RabbitMQ-Failed-to-Declare-a-Queue: Server=%s, QueueName=%s, errorMessage=%s", server, q, err.Error())
		} else {
			if conf.AutoConsume {
				Consume(conf.GroupId, q)
			}
		}
	}

}

func QueueDeclare(queueName string) error {

	_, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		logger.Errorf("RabbitMQ-Failed-to-Declare-a-Queue: Server=%s, QueueName=%s, errorMessage=%s", server, queueName, err.Error())
		return err
	}

	logger.Infof("RabbitMQ-Declare-Queue-Success: Server=%s, QueueName=%s", server, queueName)

	return nil

}

// 交换器、路由键、绑定和队列:
// -----------------------------------
// 队列通过路由键绑定到交换器, 将消息发送给交换器时, 根据每种交换器的路由规则(这里的规则就是路由键), RabbitMQ将会决定将该消息投递到哪个队列。
func Publish(queueName string, message []byte) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := channel.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		})

	if err != nil {
		logger.Errorf("RabbitMQ-Failed-to-Publish-a-Message: Server=%s, QueueName=%s, errorMessage=%s", server, queueName, err.Error())
		return
	}

	logger.Infof("RabbitMQ-Send-a-Message: [x] %s", string(message))

}

func Consume(groupId string, queueName string) {

	go func() {

		msgs, err := channel.Consume(
			queueName, // queue
			groupId,   // consumer, consumerTag
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)

		if err != nil {
			logger.Errorf("RabbitMQ-Failed-to-Consumer-a-Message: GroupId=%s, QueueName=%s, errorMessage=%s", groupId, queueName, err.Error())
			return
		}

		go func() {
			for d := range msgs {
				logger.Infof(`RabbitMQ-Received-a-message: %s`, d.Body)
			}
		}()

		logger.Infof("RabbitMQ-Waiting-for-Messages: [*] GroupId=%s, QueueName=%s", groupId, queueName)

	}()

}
