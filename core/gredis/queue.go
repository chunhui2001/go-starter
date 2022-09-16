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

	SendMessage(utils.MapOf("index", 9, "啊啊舒服的", "你好啊"))

}

func InitConsumer(sVersion string) {

	c, err := NewConsumerWithOptions(&ConsumerOptions{
		// Name:              "",
		// GroupName:         "",
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

	consumer = c

	RegisterConsumter("redisqueue:test", process)

	go func() {
		for err := range c.Errors {
			// handle errors accordingly
			logger.Errorf(fmt.Sprintf("Redis-Queue-Consumer-Handler-Message-Error: errorMessage=%+v", err))
		}
	}()

	go c.Run()

	logger.Info(fmt.Sprintf("Redis-Queue-CreateConsumer-Completed: ServerVersion=%s, %s", serverVersion, serverInfo))

}

// 发送一条消息
func SendMessage(msg map[string]interface{}) {

	err := producer.Enqueue(&Message{
		Stream: "redisqueue:test",
		Values: msg,
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Redis-Queue-SendMessage-Error: %s, ErrorMessage=%s", serverInfo, utils.ErrorToString(err)))
		return
	}

}

// 注册一个消费者
func RegisterConsumter(queueName string, process ConsumerFunc) {

	consumer.Register(queueName, process)

}

func process(msg *Message) error {
	logger.Info(fmt.Sprintf("Redis-Queue-Consumer-Processing-Message: Message=%s", utils.ToJsonString(msg)))
	return nil
}
