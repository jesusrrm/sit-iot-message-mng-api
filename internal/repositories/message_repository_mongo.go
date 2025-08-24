package repositories

import (
	"context"
	"errors"
	"log"
	"sit-iot-message-mng-api/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type messageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepository{
		collection: db.Collection("messages"), // Assuming your collection is named "messages"
	}
}

func (r *messageRepository) FindByID(ctx context.Context, id string) (*models.Message, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid message ID format")
	}

	var message models.Message
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&message)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("message not found")
		}
		return nil, err
	}
	return &message, nil
}

func (r *messageRepository) List(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	// Convert filter to bson.M
	bsonFilter := bson.M{}
	for key, value := range filter {
		if key == "_id" {
			idStr, ok := value.(string)
			if ok {
				objectID, err := primitive.ObjectIDFromHex(idStr)
				if err != nil {
					return nil, 0, errors.New("invalid ObjectID format in filter")
				}
				bsonFilter[key] = objectID
			} else {
				return nil, 0, errors.New("invalid filter value for _id")
			}
		} else {
			bsonFilter[key] = value
		}
	}

	// Sorting - default to timestamp descending for recent messages first
	sort := -1 // Default to descending for timestamp
	if sortField == "" {
		sortField = "timestamp"
	}
	if sortOrder == "ASC" {
		sort = 1
	}

	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: sortField, Value: sort}})

	// Count total documents matching the filter
	total, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		log.Printf("Error counting documents: %v", err)
		return nil, 0, err
	}

	// Find documents
	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		log.Printf("Error finding documents: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			log.Printf("Error decoding document: %v", err)
			return nil, 0, err
		}
		messages = append(messages, &message)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v", err)
		return nil, 0, err
	}

	return messages, int(total), nil
}

func (r *messageRepository) FindByTopic(ctx context.Context, topic string, limit int) ([]*models.Message, error) {
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Most recent first

	filter := bson.M{"topic": topic}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, cursor.Err()
}

func (r *messageRepository) FindByDeviceID(ctx context.Context, deviceID string, limit int) ([]*models.Message, error) {
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Most recent first

	filter := bson.M{"device_id": deviceID}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, cursor.Err()
}

func (r *messageRepository) FindByTimeRange(ctx context.Context, from, to time.Time, limit int) ([]*models.Message, error) {
	opts := options.Find()
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: "timestamp", Value: -1}}) // Most recent first

	filter := bson.M{
		"timestamp": bson.M{
			"$gte": from,
			"$lte": to,
		},
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, err
		}
		messages = append(messages, &message)
	}

	return messages, cursor.Err()
}
