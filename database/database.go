package database

import (
	"context"
	"sit-iot-message-mng-api/config"

	firebase "firebase.google.com/go/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/api/option"
)

func InitDB(cfg *config.Config) (*mongo.Database, error) {
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
