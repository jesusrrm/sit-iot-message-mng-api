package repositories

import (
	"context"

	"sit-iot-message-mng-api/internal/models"
	"time"

)

type MessageRepository interface {
	FindByID(ctx context.Context, id string) (*models.Message, error)
	List(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	FindByTopic(ctx context.Context, topic string, limit int) ([]*models.Message, error)
	FindByClientID(ctx context.Context, clientID string, limit int) ([]*models.Message, error)
	FindByTimeRange(ctx context.Context, from, to time.Time, limit int) ([]*models.Message, error)
}
