package repositories

import (
	"context"
	"errors"
	"log"
	"sit-iot-message-mng-api/internal/middleware"
	"sit-iot-message-mng-api/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	FindByID(ctx context.Context, id string) (*models.Message, error)
	Update(ctx context.Context, message *models.Message) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	ListByProjectID(ctx context.Context, projectID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
	ListByDeviceID(ctx context.Context, deviceID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error)
}

type messageRepository struct {
	collection *mongo.Collection
}

func NewMessageRepository(db *mongo.Database) MessageRepository {
	return &messageRepository{
		collection: db.Collection("messages"),
	}
}

func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return errors.New("user ID not found in context")
	}

	userEmail, ok := ctx.Value(middleware.UserEmailKey).(string)
	if !ok || userEmail == "" {
		return errors.New("user email not found in context")
	}

	// Set metadata
	message.CreatedBy = userEmail
	message.CreatedAt = time.Now()
	message.UpdatedAt = time.Now()
	message.Timestamp = time.Now()

	if message.Status == "" {
		message.Status = models.MessageStatusPending
	}

	log.Printf("Creating message: %+v", message)
	_, err := r.collection.InsertOne(ctx, message)
	return err
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

func (r *messageRepository) Update(ctx context.Context, message *models.Message) error {
	objectID, err := primitive.ObjectIDFromHex(message.ID.Hex())
	if err != nil {
		return errors.New("invalid message ID format")
	}

	message.UpdatedAt = time.Now()
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": message}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *messageRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid message ID format")
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *messageRepository) List(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
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

	sort := 1
	if sortOrder == "DESC" {
		sort = -1
	}
	opts := options.Find()
	opts.SetSkip(int64(skip))
	opts.SetLimit(int64(limit))
	opts.SetSort(bson.D{{Key: sortField, Value: sort}})

	total, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var message models.Message
		if err := cursor.Decode(&message); err != nil {
			return nil, 0, err
		}
		messages = append(messages, &message)
	}

	return messages, int(total), nil
}

func (r *messageRepository) ListByProjectID(ctx context.Context, projectID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	if filter == nil {
		filter = make(map[string]interface{})
	}
	filter["projectId"] = projectID
	return r.List(ctx, filter, sortField, sortOrder, skip, limit)
}

func (r *messageRepository) ListByDeviceID(ctx context.Context, deviceID string, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	if filter == nil {
		filter = make(map[string]interface{})
	}
	filter["deviceId"] = deviceID
	return r.List(ctx, filter, sortField, sortOrder, skip, limit)
}
