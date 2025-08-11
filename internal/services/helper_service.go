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
	"sit-iot-message-mng-api/internal/middleware"

	"time"

	"firebase.google.com/go/v4/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *messageService) VerifyToken(ctx context.Context, token string) (*auth.Token, error) {
	return s.firebaseAuth.VerifyIDToken(ctx, token)
}

// Helper function to fetch project IDs from the REST API
func (s *messageService) fetchProjectIDs(ctx context.Context) ([]primitive.ObjectID, error) {
	apiURL := s.Config.ProjectServiceApiUrl + "/api/project"

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, errors.New("failed to create HTTP request")
	}

	// Extract Authorization header from context
	tokenStr, ok := ctx.Value(middleware.TokenKey).(string)
	if !ok || tokenStr == "" {
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

// UserData represents a user with client IDs and project ID for API responses
type UserData struct {
	Username  string   `json:"username"`
	ClientIDs []string `json:"client_ids"`
	ProjectID string   `json:"project_id"`
}

// UsersResponse holds both users and their associated project IDs
type UsersResponse struct {
	Users      []UserData           `json:"users"`
	ProjectIDs []primitive.ObjectID `json:"project_ids"`
}

// Helper function to fetch all users with clients IDs from the REST API
func (s *messageService) fetchUsers(ctx context.Context) (*UsersResponse, error) {
	apiURL := s.Config.MqttServiceApiUrl + "/api/mqtt/users"

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, errors.New("failed to create HTTP request")
	}

	// Extract Authorization header from context
	tokenStr, ok := ctx.Value(middleware.TokenKey).(string)
	if !ok || tokenStr == "" {
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
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP request failed: %v", err)
		return nil, errors.New("failed to connect to MQTT API")
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
		return nil, errors.New(fmt.Sprintf("project API returned status %d: %s", resp.StatusCode, resp.Status))
	}

	// Parse the response as UserWithClients (API response format)
	var users []UserData

	// Create a new reader from the body bytes
	bodyReader := bytes.NewReader(bodyBytes)
	if err := json.NewDecoder(bodyReader).Decode(&users); err != nil {
		return nil, errors.New("failed to parse users API response")
	}

	// Extract unique project IDs from users and convert to ObjectIDs
	projectIDSet := make(map[string]bool)
	var projectIDs []primitive.ObjectID

	for i, user := range users {
		// Skip if project ID already processed
		if projectIDSet[user.ProjectID] {
			continue
		}
		projectIDSet[user.ProjectID] = true

		objectID, err := primitive.ObjectIDFromHex(user.ProjectID)
		if err != nil {
			log.Printf("Invalid ObjectID format for user %d project ID (%s): %v", i+1, user.ProjectID, err)
			return nil, errors.New(fmt.Sprintf("invalid project ID format: %s", user.ProjectID))
		}
		projectIDs = append(projectIDs, objectID)
	}

	return &UsersResponse{
		Users:      users,
		ProjectIDs: projectIDs,
	}, nil
}
