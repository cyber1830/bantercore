package platform

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/sahilverma/muze-go-backend/discussions"
)

type NoopCache struct {
	mu   sync.RWMutex
	data map[string][]discussions.Discussion
}

func NewNoopCache() *NoopCache {
	return &NoopCache{data: map[string][]discussions.Discussion{}}
}

func (c *NoopCache) Get(_ context.Context, key string) ([]discussions.Discussion, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

func (c *NoopCache) Set(_ context.Context, key string, value []discussions.Discussion, _ time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}

func (c *NoopCache) Delete(_ context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

type LogPublisher struct {
	logger *slog.Logger
}

func NewLogPublisher(logger *slog.Logger) *LogPublisher {
	return &LogPublisher{logger: logger}
}

func (p *LogPublisher) Publish(_ context.Context, eventType string, value any) error {
	p.logger.Info("event_published", "type", eventType, "payload", value)
	return nil
}
