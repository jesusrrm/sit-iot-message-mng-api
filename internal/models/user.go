package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserWithClientIDs represents a user with their associated client IDs for repository operations
type UserWithClientIDs struct {
	Username  string             `bson:"username" json:"username"`
	ClientIDs []string           `bson:"client_ids" json:"client_ids"`
	ProjectID primitive.ObjectID `bson:"project_id" json:"project_id"`
}

// UserWithClients represents a user with their associated client IDs for API responses
type UserWithClients struct {
	Username  string   `json:"username"`
	ClientIDs []string `json:"client_ids"`
	ProjectID string   `json:"project_id"`
}