package services

import (
	"context"
	"errors"

	"log"
	"sit-iot-message-mng-api/config"
	"sit-iot-message-mng-api/internal/middleware"
	"sit-iot-message-mng-api/internal/models"
	"sit-iot-message-mng-api/internal/repositories"

	"firebase.google.com/go/v4/auth"
)


type messageService struct {
	messageRepo  repositories.MessageRepository
	firebaseAuth *auth.Client
	Config      *config.Config

}

func NewMessageService(messageRepo repositories.MessageRepository, firebaseAuth *auth.Client) MessageService {
	return &messageService{
		messageRepo:  messageRepo,
		firebaseAuth: firebaseAuth,
		Config:      &config.Config{},
	}
}



func (s *messageService) GetMessageByID(ctx context.Context, id string) (*models.Message, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("user ID not found in context")
	}

	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return nil, errors.New("user email not found in context3")
	}

	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// TODO: Add authorization check to ensure user can access this message
	// For now, we'll return the message

	return message, nil
}


func (s *messageService) ListMessages(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return nil, 0, errors.New("user email not found in context4")
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

	return []*models.Message{}, 0, nil //s.messageRepo.ListByProjectID(ctx, projectID, filter, sortField, sortOrder, skip, limit)
}

func (s *messageService) ListMessagesByDeviceID(ctx context.Context, deviceID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return nil, 0, errors.New("user email not found in context2")
	}

	// TODO: Add authorization check to ensure user can access messages from this device
		// Make a REST call to fetch project IDs to ensure user has access
	projectIDs, err := s.fetchProjectIDs(ctx)
	if err != nil {
		return nil, 0, err
	}

	if projectIDs != nil {
		log.Printf("Token verification failed: %v", err)
		return nil, 0, err
	}

	return []*models.Message{}, 0, nil // s.messageRepo.ListByDeviceID(ctx, deviceID, filter, sortField, sortOrder, skip, limit)
}