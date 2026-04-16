package APIs

import (
	"API/internal/errs"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"API/internal/config+conn"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ARepository interface {
	Create(ctx context.Context, userID uuid.UUID, title string) error

	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*API, error)
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
	Update(ctx context.Context, id, userID uuid.UUID, title *string, completed *bool) error
	GetAll(ctx context.Context, userID uuid.UUID, filter QueryFilter) ([]API, error)
}

type APService struct {
	repo  ARepository
	redis *redis.Client
}

func NewAPService(repo ARepository, redis *redis.Client) *APService {
	return &APService{
		repo:  repo,
		redis: redis,
	}
}
func buildAPICacheKey(userID uuid.UUID, filter QueryFilterFromHandler) string {
	completed := "nil"
	if filter.Completed != nil {
		completed = strconv.FormatBool(*filter.Completed)
	}

	title := ""
	if filter.Title != nil {
		title = *filter.Title
	}

	return fmt.Sprintf(
		"Api:%s:limit:%d:offset:%d:completed:%s:title:%s",
		userID.String(),
		*filter.Limit,
		*filter.Offset,
		completed,
		title,
	)
}

func (s *APService) Create(ctx context.Context, userID uuid.UUID, title string) error {
	if userID == uuid.Nil {
		return fmt.Errorf("userID should not be nil %w", errs.ErrBadRequest)
	}
	if title == "" {
		return fmt.Errorf("title should not be empty %w", errs.ErrBadRequest)
	}

	err := s.repo.Create(ctx, userID, title)
	if err != nil {
		return fmt.Errorf("repo.Create: %w", err)
	}
	if err := config_conn.InvalidateCache(userID, s.redis); err != nil {
		log.Printf("failed to invalidate cache: %v", err)
	}
	return nil
}
func (s *APService) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return fmt.Errorf("userID should not be nil: %w", errs.ErrBadRequest)
	}
	if id == uuid.Nil {
		return fmt.Errorf("id should not be nil: %w", errs.ErrBadRequest)
	}

	if err := s.repo.Delete(ctx, userID, id); err != nil {
		return fmt.Errorf("error delete: %w", err)
	}

	if err := config_conn.InvalidateCache(userID, s.redis); err != nil {
		log.Printf("failed to invalidate cache: %v", err)
	}

	return nil
}
func (s *APService) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*API, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("userID should not be nil %w", errs.ErrBadRequest)
	}
	if id == uuid.Nil {
		return nil, fmt.Errorf("id should not be nil %w", errs.ErrBadRequest)
	}
	api, err := s.repo.GetByID(ctx, userID, id)
	if err != nil {
		return nil, fmt.Errorf("error GetByID: %w", err)
	}
	return api, nil
}

func (s *APService) GetAll(ctx context.Context, userID uuid.UUID, filter QueryFilterFromHandler) ([]API, error) {
	if userID == uuid.Nil {
		return nil, fmt.Errorf("userID should not be nil: %w", errs.ErrBadRequest)
	}

	if filter.Limit == nil || *filter.Limit < 2 || *filter.Limit > 100 {
		defaultLimit := 10
		filter.Limit = &defaultLimit
	}

	if filter.Offset == nil || *filter.Offset < 0 {
		defaultOffset := 0
		filter.Offset = &defaultOffset
	}

	key := buildAPICacheKey(userID, filter)
	cached, err := s.redis.Get(ctx, key).Result()

	if err == nil {
		var apis []API
		err = json.Unmarshal([]byte(cached), &apis)
		if err == nil {
			return apis, nil
		}
		log.Printf("failed to unmarshal cached APIs: %v", err)
	} else if !errors.Is(err, redis.Nil) {
		log.Printf("failed to cached: %v", err)
	}

	apis, err := s.repo.GetAll(ctx, userID, QueryFilter{
		completed: filter.Completed,
		title:     filter.Title,
		limit:     filter.Limit,
		offset:    filter.Offset,
	})
	if err != nil {
		return nil, fmt.Errorf("get all apis: %w", err)
	}
	date, err := json.Marshal(apis)
	if err != nil {
		log.Printf("failed to marshal APIs: %v", err)
		return apis, nil
	}
	if err := s.redis.Set(ctx, key, date, time.Minute).Err(); err != nil {
		log.Printf("failed to set cached APIs: %v", err)
	}

	return apis, nil
}

func (s *APService) Update(ctx context.Context, id uuid.UUID, userID uuid.UUID, title *string, completed *bool) error {
	if userID == uuid.Nil {
		return fmt.Errorf("userID should not be nil: %w", errs.ErrBadRequest)
	}
	if id == uuid.Nil {
		return fmt.Errorf("id should not be nil: %w", errs.ErrBadRequest)
	}
	if title == nil && completed == nil {
		return fmt.Errorf("title should not be nil: %w", errs.ErrBadRequest)
	}

	if err := s.repo.Update(ctx, userID, id, title, completed); err != nil {
		return fmt.Errorf("repo.Update: %w", err)
	}
	if err := config_conn.InvalidateCache(userID, s.redis); err != nil {
		log.Printf("failed to invalidate cache: %v", err)
	}
	return nil

}
