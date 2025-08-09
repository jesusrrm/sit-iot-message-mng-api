package services

import (
	"context"
	"errors"
	"log"
	"sit-iot-message-mng-api/internal/middleware"
	"sit-iot-message-mng-api/internal/models"
	"sit-iot-message-mng-api/internal/repositories"

	"firebase.google.com/go/v4/auth"
)

type MessageService interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessageByID(ctx context.Context, id string) (*models.Message, error)
	UpdateMessage(ctx context.Context, message *models.Message) error
	DeleteMessage(ctx context.Context, id string) error
	ListMessages(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	ListMessagesByProjectID(ctx context.Context, projectID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	ListMessagesByDeviceID(ctx context.Context, deviceID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	VerifyToken(ctx context.Context, token string) (*auth.Token, error)
}

type messageService struct {
	messageRepo  repositories.MessageRepository
	firebaseAuth *auth.Client
}

func NewMessageService(messageRepo repositories.MessageRepository, firebaseAuth *auth.Client) MessageService {
	return &messageService{
		messageRepo:  messageRepo,
		firebaseAuth: firebaseAuth,
	}
}

func (s *messageService) CreateMessage(ctx context.Context, message *models.Message) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errors.New("user ID not found in context")
	}

	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return errors.New("user email not found in context")
	}

	log.Printf("Service: Creating message for user %s (%s)", userID, userEmail)

	// TODO: Add project membership validation here
	// For now, we'll just create the message

	return s.messageRepo.Create(ctx, message)
}

func (s *messageService) GetMessageByID(ctx context.Context, id string) (*models.Message, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("user ID not found in context")
	}

	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return nil, errors.New("user email not found in context")
	}

	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Add authorization check to ensure user can access this message
	// For now, we'll return the message

	return message, nil
}

func (s *messageService) UpdateMessage(ctx context.Context, message *models.Message) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errors.New("user ID not found in context")
	}

	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return errors.New("user email not found in context")
	}

	// TODO: Add authorization check to ensure user can update this message

	return s.messageRepo.Update(ctx, message)
}

func (s *messageService) DeleteMessage(ctx context.Context, id string) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errors.New("user ID not found in context")
	}

	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return errors.New("user email not found in context")
	}

	// TODO: Add authorization check to ensure user can delete this message

	return s.messageRepo.Delete(ctx, id)
}

func (s *messageService) ListMessages(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return nil, 0, errors.New("user email not found in context")
	}

	// TODO: Add filter to only show messages from projects where user is a member
	// For now, we'll return all messages that the user created

	if filter == nil {
		filter = make(map[string]interface{})
	}
	filter["createdBy"] = userEmail

	return s.messageRepo.List(ctx, filter, sortField, sortOrder, skip, limit)
}

func (s *messageService) ListMessagesByProjectID(ctx context.Context, projectID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return nil, 0, errors.New("user email not found in context")
	}

	// TODO: Add authorization check to ensure user is a member of the project

	return s.messageRepo.ListByProjectID(ctx, projectID, filter, sortField, sortOrder, skip, limit)
}

func (s *messageService) ListMessagesByDeviceID(ctx context.Context, deviceID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return nil, 0, errors.New("user email not found in context")
	}

	// TODO: Add authorization check to ensure user can access messages from this device

	return s.messageRepo.ListByDeviceID(ctx, deviceID, filter, sortField, sortOrder, skip, limit)
}

func (s *messageService) VerifyToken(ctx context.Context, token string) (*auth.Token, error) {
	return s.firebaseAuth.VerifyIDToken(ctx, token)
}
