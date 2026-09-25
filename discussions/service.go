package discussions

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Discussion struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"authorId"`
	Topic     string    `json:"topic"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type Store interface {
	Create(context.Context, Discussion) error
	List(context.Context, int) ([]Discussion, error)
}

type Cache interface {
	Get(context.Context, string) ([]Discussion, bool)
	Set(context.Context, string, []Discussion, time.Duration)
}

type Publisher interface {
	Publish(context.Context, string, any) error
}

type Service struct {
	store     Store
	cache     Cache
	publisher Publisher
}

func NewService(store Store, cache Cache, publisher Publisher) *Service {
	return &Service{store: store, cache: cache, publisher: publisher}
}

func (s *Service) Create(ctx context.Context, author, topic, body string) (Discussion, error) {
	topic, body = strings.TrimSpace(topic), strings.TrimSpace(body)
	if author == "" || topic == "" || body == "" {
		return Discussion{}, errors.New("author, topic, and body are required")
	}

	discussion := Discussion{
		ID:        uuid.NewString(),
		AuthorID:  author,
		Topic:     topic,
		Body:      body,
		CreatedAt: time.Now().UTC(),
	}
	if err := s.store.Create(ctx, discussion); err != nil {
		return Discussion{}, err
	}
	_ = s.publisher.Publish(ctx, "discussion.created", discussion)
	return discussion, nil
}

func (s *Service) List(ctx context.Context, limit int) ([]Discussion, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if cached, ok := s.cache.Get(ctx, "discussions"); ok {
		return cached, nil
	}

	rows, err := s.store.List(ctx, limit)
	if err == nil {
		s.cache.Set(ctx, "discussions", rows, 30*time.Second)
	}
	return rows, err
}
