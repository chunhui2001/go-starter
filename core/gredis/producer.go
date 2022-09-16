package gredis

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
)

var redisVersionRE = regexp.MustCompile(`redis_version:(.+)`)

// Message constitutes a message that will be enqueued and dequeued from Redis.
// When enqueuing, it's recommended to leave ID empty and let Redis generate it,
// unless you know what you're doing.
type Message struct {
	ID     string
	Stream string
	Values map[string]interface{}
}

// ProducerOptions provide options to configure the Producer.
type ProducerOptions struct {
	// StreamMaxLength sets the MAXLEN option when calling XADD. This creates a
	// capped stream to prevent the stream from taking up memory indefinitely.
	// It's important to note though that this isn't the maximum number of
	// _completed_ messages, but the maximum number of _total_ messages. This
	// means that if all consumers are down, but producers are still enqueuing,
	// and the maximum is reached, unprocessed message will start to be dropped.
	// So ideally, you'll set this number to be as high as you can makee it.
	// More info here: https://redis.io/commands/xadd#capped-streams.
	StreamMaxLength int64
	// ApproximateMaxLength determines whether to use the ~ with the MAXLEN
	// option. This allows the stream trimming to done in a more efficient
	// manner. More info here: https://redis.io/commands/xadd#capped-streams.
	ApproximateMaxLength bool
	// RedisClient supersedes the RedisOptions field, and allows you to inject
	// an already-made Redis Client for use in the consumer. This may be either
	// the standard client or a cluster client.
	RedisClient redis.UniversalClient
	// RedisOptions allows you to configure the underlying Redis connection.
	// More info here:
	// https://pkg.go.dev/github.com/go-redis/redis/v7?tab=doc#Options.
	//
	// This field is used if RedisClient field is nil.
	RedisOptions *redis.Options
}

// Producer adds a convenient wrapper around enqueuing messages that will be
// processed later by a Consumer.
type Producer struct {
	options *ProducerOptions
	redis   redis.UniversalClient
}

var defaultProducerOptions = &ProducerOptions{
	StreamMaxLength:      1000,
	ApproximateMaxLength: true,
}

// NewProducer uses a default set of options to create a Producer. It sets
// StreamMaxLength to 1000 and ApproximateMaxLength to true. In most production
// environments, you'll want to use NewProducerWithOptions.
func NewProducer() (*Producer, error) {
	return NewProducerWithOptions(defaultProducerOptions)
}

// NewProducerWithOptions creates a Producer using custom ProducerOptions.
func NewProducerWithOptions(options *ProducerOptions) (*Producer, error) {

	var r redis.UniversalClient = options.RedisClient

	if err := redisPreflightChecks(r); err != nil {
		return nil, err
	}

	return &Producer{
		options: options,
		redis:   r,
	}, nil

}

// redisPreflightChecks makes sure the Redis instance backing the *redis.Client
// offers the functionality we need. Specifically, it also that it can connect
// to the actual instance and that the instance supports Redis streams (i.e.
// it's at least v5).
func redisPreflightChecks(client redis.UniversalClient) error {

	info, err := client.Info(context.Background(), "server").Result()

	if err != nil {
		return err
	}

	match := redisVersionRE.FindAllStringSubmatch(info, -1)

	if len(match) < 1 {
		return fmt.Errorf("could not extract redis version")
	}

	version := strings.TrimSpace(match[0][1])
	parts := strings.Split(version, ".")
	major, err := strconv.Atoi(parts[0])

	if err != nil {
		return err
	}

	if major < 5 {
		return fmt.Errorf("redis streams are not supported in version %q", version)
	}

	return nil
}

// Enqueue takes in a pointer to Message and enqueues it into the stream set at
// msg.Stream. While you can set msg.ID, unless you know what you're doing, you
// should let Redis auto-generate the ID. If an ID is auto-generated, it will be
// set on msg.ID for your reference. msg.Values is also required.
func (p *Producer) Enqueue(msg *Message) error {
	args := &redis.XAddArgs{
		ID:     msg.ID,
		Stream: msg.Stream,
		Values: msg.Values,
	}
	if p.options.ApproximateMaxLength {
		args.MaxLenApprox = p.options.StreamMaxLength
	} else {
		args.MaxLen = p.options.StreamMaxLength
	}
	id, err := p.redis.XAdd(context.Background(), args).Result()
	if err != nil {
		return err
	}
	msg.ID = id
	return nil
}
