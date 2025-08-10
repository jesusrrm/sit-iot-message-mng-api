package services

import (
	"context"
	"sit-iot-message-mng-api/internal/models"

	"firebase.google.com/go/v4/auth"
)

type MessageService interface {
	GetMessageByID(ctx context.Context, id string) (*models.Message, error)
	ListMessages(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	ListMessagesByProjectID(ctx context.Context, projectID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	ListMessagesByDeviceID(ctx context.Context, deviceID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	VerifyToken(ctx context.Context, token string) (*auth.Token, error)
}
