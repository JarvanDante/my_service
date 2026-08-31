// Package mq 站点侧 Kafka：投递换脸/去衣任务、消费结果。
package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	kafkago "github.com/segmentio/kafka-go"

	"github.com/JarvanDante/my_service/internal/shared/aiprotocol"
)

// Bus 换脸 jobs/results。不要复用转码 topic。
type Bus struct {
	enabled      bool
	brokers      []string
	topicJobs    string
	topicResults string
	group        string
	writer       *kafkago.Writer
}

var (
	mu       sync.Mutex
	defaultB *Bus
)

func Init(ctx context.Context) *Bus {
	mu.Lock()
	defer mu.Unlock()
	if defaultB != nil {
		return defaultB
	}
	enabled := g.Cfg().MustGet(ctx, "kafka.enabled", true).Bool()
	brokers := g.Cfg().MustGet(ctx, "kafka.brokers").Strings()
	if len(brokers) == 0 {
		brokers = []string{"127.0.0.1:9092"}
	}
	jobs := g.Cfg().MustGet(ctx, "kafka.faceswap_job_topic", aiprotocol.TopicJobs).String()
	results := g.Cfg().MustGet(ctx, "kafka.faceswap_result_topic", aiprotocol.TopicResults).String()
	group := g.Cfg().MustGet(ctx, "kafka.groupId", "my_service").String()
	if group == "" {
		group = "my_service"
	}
	b := &Bus{
		enabled:      enabled,
		brokers:      brokers,
		topicJobs:    jobs,
		topicResults: results,
		group:        group,
	}
	if enabled {
		b.writer = &kafkago.Writer{
			Addr:                   kafkago.TCP(brokers...),
			Topic:                  jobs,
			Balancer:               &kafkago.Hash{},
			RequiredAcks:           kafkago.RequireOne,
			Async:                  false,
			AllowAutoTopicCreation: true,
		}
		ensureCtx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		if err := b.ensureTopics(ensureCtx); err != nil {
			g.Log().Warningf(ctx, "kafka: ensure faceswap topics: %v", err)
		}
		cancel()
	}
	defaultB = b
	return b
}

func Default() *Bus {
	if defaultB != nil {
		return defaultB
	}
	return Init(context.Background())
}

func (b *Bus) Enabled() bool {
	return b != nil && b.enabled
}

func (b *Bus) Close() error {
	if b == nil || b.writer == nil {
		return nil
	}
	return b.writer.Close()
}

func (b *Bus) PublishJob(ctx context.Context, msg aiprotocol.JobMessage) error {
	if b == nil || !b.enabled {
		return fmt.Errorf("kafka disabled")
	}
	if b.writer == nil {
		return fmt.Errorf("kafka writer not initialized")
	}
	if msg.JobID == "" {
		return fmt.Errorf("job_id required")
	}
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return b.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(msg.JobID),
		Value: body,
		Time:  time.Now(),
	})
}

func (b *Bus) ConsumeResults(ctx context.Context, handler func(context.Context, aiprotocol.ResultMessage) error) error {
	if b == nil || !b.enabled {
		g.Log().Warning(ctx, "kafka disabled, skip faceswap result consumer")
		<-ctx.Done()
		return ctx.Err()
	}
	const minBackoff, maxBackoff = 2 * time.Second, 30 * time.Second
	backoff := minBackoff
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := b.consumeResultsOnce(ctx, handler)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		g.Log().Warningf(ctx, "kafka: faceswap result consumer exited: %v; reconnect in %s", err, backoff)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		if backoff < maxBackoff {
			backoff *= 2
			if backoff > maxBackoff {
				backoff = maxBackoff
			}
		}
	}
}

func (b *Bus) consumeResultsOnce(ctx context.Context, handler func(context.Context, aiprotocol.ResultMessage) error) error {
	r := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        b.brokers,
		GroupID:        b.group,
		Topic:          b.topicResults,
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
		StartOffset:    kafkago.FirstOffset,
	})
	defer r.Close()

	g.Log().Infof(ctx, "kafka: consuming faceswap results topic=%s group=%s brokers=%s",
		b.topicResults, b.group, strings.Join(b.brokers, ","))

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("kafka fetch: %w", err)
		}
		var res aiprotocol.ResultMessage
		if err := json.Unmarshal(m.Value, &res); err != nil {
			g.Log().Warningf(ctx, "kafka: invalid faceswap result json: %v; skip", err)
			_ = r.CommitMessages(ctx, m)
			continue
		}
		if err := handler(ctx, res); err != nil {
			g.Log().Warningf(ctx, "kafka: faceswap result handle failed job_id=%s: %v", res.JobID, err)
			continue
		}
		if err := r.CommitMessages(ctx, m); err != nil {
			g.Log().Warningf(ctx, "kafka: commit failed: %v", err)
		}
	}
}

func (b *Bus) ensureTopics(ctx context.Context) error {
	if len(b.brokers) == 0 {
		return fmt.Errorf("kafka brokers empty")
	}
	for _, topic := range []string{b.topicJobs, b.topicResults} {
		if err := createTopic(ctx, b.brokers[0], topic); err != nil {
			return fmt.Errorf("ensure topic %s: %w", topic, err)
		}
	}
	return nil
}

func createTopic(ctx context.Context, broker, topic string) error {
	conn, err := kafkago.DialContext(ctx, "tcp", broker)
	if err != nil {
		return err
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	ctrlAddr := net.JoinHostPort(controller.Host, fmt.Sprintf("%d", controller.Port))
	ctrl, err := kafkago.DialContext(ctx, "tcp", ctrlAddr)
	if err != nil {
		return err
	}
	defer ctrl.Close()
	err = ctrl.CreateTopics(kafkago.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		low := strings.ToLower(err.Error())
		if strings.Contains(low, "already exists") {
			return nil
		}
		return err
	}
	return nil
}
