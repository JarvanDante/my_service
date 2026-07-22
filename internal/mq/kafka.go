// Package mq Kafka 生产/消费封装(占位)。
// 建议引入 github.com/segmentio/kafka-go 或 franz-go, 生产者/消费者在此收口。
package mq

import "context"

type Producer interface {
	Send(ctx context.Context, topic string, key, value []byte) error
}

// TODO: 实现基于 kafka-go 的 Producer, 从配置读取 brokers。
