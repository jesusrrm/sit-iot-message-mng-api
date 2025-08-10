package database

import (
	"context"
	"sit-iot-message-mng-api/config"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

// DatabaseClients holds both MongoDB and Firestore clients
type DatabaseClients struct {
	MongoDB   *mongo.Database
	Firestore *firestore.Client
}

// InitDatabases initializes the appropriate database based on configuration
func InitDatabases(cfg *config.Config) (*DatabaseClients, error) {
	clients := &DatabaseClients{}

	switch cfg.DatabaseProvider {
	case "mongo", "mongodb":
		mongoClient, err := InitMongoDB(cfg)
		if err != nil {
			return nil, err
		}
		clients.MongoDB = mongoClient

	case "firestore":
		firestoreClient, err := InitFirestore(cfg.FirebaseCredentialsPath)
		if err != nil {
			return nil, err
		}
		clients.Firestore = firestoreClient

	default:
		// Initialize both for backwards compatibility or if not specified
		mongoClient, err := InitMongoDB(cfg)
		if err == nil {
			clients.MongoDB = mongoClient
		}

		firestoreClient, err := InitFirestore(cfg.FirebaseCredentialsPath)
		if err == nil {
			clients.Firestore = firestoreClient
		}
	}

	return clients, nil
}

func InitMongoDB(cfg *config.Config) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(cfg.DatabaseURL)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	// Extract database name from URL (simplified)
	db := client.Database(cfg.DBName)
	return db, nil
}

func InitFirestore(credentialsPath string) (*firestore.Client, error) {
	ctx := context.Background()
	var client *firestore.Client
	var err error

	if credentialsPath != "" {
		client, err = firestore.NewClient(ctx, "", option.WithCredentialsFile(credentialsPath))
	} else {
		// Use default credentials (for GCP environments)
		client, err = firestore.NewClient(ctx, "")
	}

	if err != nil {
		return nil, err
	}

	return client, nil
}

func InitFirebase(credentialsPath string) (*firebase.App, error) {
	var opt option.ClientOption
	if credentialsPath != "" {
		opt = option.WithCredentialsFile(credentialsPath)
	} else {
		// Use default credentials
		opt = option.WithCredentialsFile("")
	}

	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return app, nil
}
