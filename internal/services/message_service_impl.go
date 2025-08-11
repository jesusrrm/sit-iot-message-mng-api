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
	Config       *config.Config
}

func NewMessageService(messageRepo repositories.MessageRepository, firebaseAuth *auth.Client, cfg *config.Config) MessageService {
	return &messageService{
		messageRepo:  messageRepo,
		firebaseAuth: firebaseAuth,
		Config:       cfg,
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
	// userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	// if !ok || userEmail == "" {
	// 	return nil, 0, errors.New("user email not found in context4")
	// }

	// TODO: Add filter to only show messages from projects where user is a member
	// For now, we'll return all messages that the user created
	usersResponse, err := s.fetchUsers(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Log the fetched users for debugging
	log.Printf("Fetched %d users from %d projects", len(usersResponse.Users), len(usersResponse.ProjectIDs))
	log.Printf("%+v", usersResponse)

	// Collect all client IDs from all users
	var allClientIDs []string
	for _, user := range usersResponse.Users {
		allClientIDs = append(allClientIDs, user.ClientIDs...)
	}

	log.Printf("Collected %d client IDs: %v", len(allClientIDs), allClientIDs)

	// Return error if no client IDs found
	if len(allClientIDs) == 0 {
		return nil, 0, errors.New("no client IDs found for the user - access denied")
	}

	if filter == nil {
		filter = make(map[string]interface{})
	}

	// Filter messages by client IDs instead of user email
	filter["client_id"] = map[string]interface{}{
		"$in": allClientIDs,
	}

	return s.messageRepo.List(ctx, filter, sortField, sortOrder, skip, limit)
}

func (s *messageService) ListMessagesByDeviceID(ctx context.Context, deviceID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return nil, 0, errors.New("user email not found in context2")
	}

	// TODO: Add authorization check to ensure user can access messages from this device
	// Make a REST call to fetch project IDs to ensure user has access
	usersResponse, err := s.fetchUsers(ctx)
	if err != nil {
		return nil, 0, err
	}

	// Validate that user has access to client IDs (similar validation as in ListMessages)
	var allClientIDs []string
	for _, user := range usersResponse.Users {
		allClientIDs = append(allClientIDs, user.ClientIDs...)
	}

	// Check if the requested deviceID (clientID) is in the user's allowed client IDs
	var hasAccess bool
	for _, clientID := range allClientIDs {
		if clientID == deviceID {
			hasAccess = true
			break
		}
	}

	if !hasAccess {
		return nil, 0, errors.New("access denied: device not found in user's allowed client IDs")
	}

	messages, err := s.messageRepo.FindByClientID(ctx, deviceID, limit)
	if err != nil {
		return nil, 0, err
	}
	return messages, len(messages), nil
}
