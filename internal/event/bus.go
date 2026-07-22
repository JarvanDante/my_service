// Package event 进程内模块间事件总线, 解耦跨模块的异步调用。
// 需要持久化/跨进程时改用 internal/mq(Kafka)。
package event

import (
	"context"
	"sync"
)

type Handler func(ctx context.Context, payload any)

type Bus struct {
	mu       sync.RWMutex
	handlers map[string][]Handler
}

var Default = &Bus{handlers: make(map[string][]Handler)}

func (b *Bus) Subscribe(topic string, h Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], h)
}

func (b *Bus) Publish(ctx context.Context, topic string, payload any) {
	b.mu.RLock()
	hs := b.handlers[topic]
	b.mu.RUnlock()
	for _, h := range hs {
		h(ctx, payload)
	}
}
