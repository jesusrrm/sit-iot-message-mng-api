package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sit-iot-message-mng-api/config"
	"sit-iot-message-mng-api/internal/middleware"
	"sit-iot-message-mng-api/internal/models"
	"sit-iot-message-mng-api/internal/repositories"
	"time"

	"firebase.google.com/go/v4/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Config      *config.Config

}

func NewMessageService(messageRepo repositories.MessageRepository, firebaseAuth *auth.Client) MessageService {
	return &messageService{
		messageRepo:  messageRepo,
		firebaseAuth: firebaseAuth,
		Config:      &config.Config{},
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
	// userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	// if !ok || userEmail == "" {
	// 	return nil, 0, errors.New("user email not found in context4")
	// }

	// TODO: Add filter to only show messages from projects where user is a member
	// For now, we'll return all messages that the user created

	// if filter == nil {
	// 	filter = make(map[string]interface{})
	// }
	// filter["createdBy"] = userEmail

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

	return s.messageRepo.ListByDeviceID(ctx, deviceID, filter, sortField, sortOrder, skip, limit)
}

func (s *messageService) VerifyToken(ctx context.Context, token string) (*auth.Token, error) {
	return s.firebaseAuth.VerifyIDToken(ctx, token)
}

// Helper function to fetch project IDs from the REST API
func (s *messageService) fetchProjectIDs(ctx context.Context) ([]primitive.ObjectID, error) {
	// Log the API URL being called
	apiURL := s.Config.ProjectServiceApiUrl + "/api/project"

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, errors.New("failed to create HTTP request")
	}

	// Extract Authorization header from context
	tokenStr, ok := ctx.Value(middleware.TokenKey).(string)
	if !ok || tokenStr == "" {
		log.Printf("Authorization token not found in context")
		return nil, errors.New("authorization token is missing")
	}

	// Add headers
	req.Header.Set("accept", "*/*")
	req.Header.Set("authorization", "Bearer "+tokenStr)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", "sit-iot-mqtt-service")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second, // Add timeout
	}

	// Make the request
	log.Printf("Making HTTP request to project API...")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, errors.New("failed to connect to project API")
	}
	defer resp.Body.Close()

	// Read response body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return nil, errors.New("failed to read API response")
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		log.Printf("Project API returned non-200 status: %d, body: %s", resp.StatusCode, string(bodyBytes))
		return nil, errors.New(fmt.Sprintf("project API returned status %d: %s", resp.StatusCode, resp.Status))
	}

	// Parse the response
	var projects []struct {
		ID string `json:"id"`
	}

	// Create a new reader from the body bytes
	bodyReader := bytes.NewReader(bodyBytes)
	if err := json.NewDecoder(bodyReader).Decode(&projects); err != nil {
		log.Printf("Failed to decode JSON response: %v", err)
		log.Printf("Response body was: %s", string(bodyBytes))
		return nil, errors.New("failed to parse project API response")
	}

	// Convert project IDs to ObjectIDs
	var projectIDs []primitive.ObjectID
	for i, project := range projects {

		objectID, err := primitive.ObjectIDFromHex(project.ID)
		if err != nil {
			log.Printf("Invalid ObjectID format for project %d (ID: %s): %v", i+1, project.ID, err)
			return nil, errors.New(fmt.Sprintf("invalid project ID format: %s", project.ID))
		}
		projectIDs = append(projectIDs, objectID)
	}

	return projectIDs, nil
}