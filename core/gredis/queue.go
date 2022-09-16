package gredis

import (
	"fmt"
	"time"

	"github.com/chunhui2001/go-starter/core/utils"
)

var (
	producer      *Producer
	consumer      *Consumer
	serverInfo    string
	serverVersion string
)

func InitProducer(sVersion string) {

	serverVersion = sVersion

	p, err := NewProducerWithOptions(&ProducerOptions{
		StreamMaxLength:      10000,
		ApproximateMaxLength: true,
		RedisClient:          universalClient,
	})

	serverInfo = fmt.Sprintf("Mode=%s, ServerAddrs=%s", conf.Mode.String(), conf.ServerAddrs())

	if err != nil {
		logger.Error(fmt.Sprintf("Redis-Queue-CreateProducer-Error: ServerVersion=%s, %s, errorMessage=%s", serverVersion, serverInfo, utils.ErrorToString(err)))
		return
	}

	producer = p

	logger.Info(fmt.Sprintf("Redis-Queue-CreateProducer-Completed: ServerVersion=%s, %s", serverVersion, serverInfo))

}

func CreateConsumer(queueName string, groupName string, process ConsumerFunc) {

	c, err := NewConsumerWithOptions(&ConsumerOptions{
		// Name:              "",
		GroupName:         groupName,
		VisibilityTimeout: 60 * time.Second,
		BlockingTimeout:   5 * time.Second,
		ReclaimInterval:   1 * time.Second,
		BufferSize:        100,
		Concurrency:       10,
		RedisClient:       universalClient,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Redis-Queue-CreateConsumer-Error: ServerVersion=%s, %s, errorMessage=%s", serverVersion, serverInfo, utils.ErrorToString(err)))
		return
	}

	c.Register(queueName, process)

	go func() {
		for err := range c.Errors {
			logger.Errorf(fmt.Sprintf("Redis-Queue-Consumer-Handler-Message-Error: errorMessage=%+v", err))
		}
	}()

	go c.Run()

	logger.Info(fmt.Sprintf("Redis-Queue-CreateConsumer-Completed: ServerVersion=%s, %s", serverVersion, serverInfo))

}

func SendMessage(queueName string, msg map[string]interface{}) bool {

	err := producer.Enqueue(&Message{
		Stream: queueName,
		Values: msg,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Redis-Queue-SendMessage-Error: %s, ErrorMessage=%s", serverInfo, utils.ErrorToString(err)))
		return false
	}

	return true

}
