package repositories

import (
	"errors"
	"sit-iot-message-mng-api/config"

	"cloud.google.com/go/firestore"
	"go.mongodb.org/mongo-driver/mongo"
)

// RepositoryFactory provides a way to create repositories based on configuration
type RepositoryFactory struct {
	config *config.Config
}

// NewRepositoryFactory creates a new repository factory
func NewRepositoryFactory(cfg *config.Config) *RepositoryFactory {
	return &RepositoryFactory{
		config: cfg,
	}
}

// CreateMessageRepository creates a message repository based on the configured database provider
func (f *RepositoryFactory) CreateMessageRepository(mongoClient *mongo.Database, firestoreClient *firestore.Client) (MessageRepository, error) {
	switch f.config.DatabaseProvider {
	case "mongo", "mongodb":
		if mongoClient == nil {
			return nil, errors.New("MongoDB client is required when using MongoDB provider")
		}
		return NewMessageRepository(mongoClient), nil
	case "firestore":
		if firestoreClient == nil {
			return nil, errors.New("Firestore client is required when using Firestore provider")
		}
		return NewFirestoreMessageRepository(firestoreClient), nil
	default:
		return nil, errors.New("unsupported database provider: " + f.config.DatabaseProvider)
	}
}

// GetDatabaseProvider returns the configured database provider
func (f *RepositoryFactory) GetDatabaseProvider() string {
	return f.config.DatabaseProvider
}
