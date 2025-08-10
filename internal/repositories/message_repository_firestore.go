package repositories

import (
	"context"
	"errors"
	"log"
    "sit-iot-message-mng-api/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type firestoreMessageRepository struct {
	client     *firestore.Client
	collection string
}

func NewFirestoreMessageRepository(client *firestore.Client) MessageRepository {
	return &firestoreMessageRepository{
		client:     client,
		collection: "messages",
	}
}


func (r *firestoreMessageRepository) FindByID(ctx context.Context, id string) (*models.Message, error) {
	if id == "" {
		return nil, errors.New("invalid message ID format")
	}

	doc, err := r.client.Collection(r.collection).Doc(id).Get(ctx)
	if err != nil {
		if doc != nil && !doc.Exists() {
			return nil, errors.New("message not found")
		}
		return nil, err
	}

	var message models.Message
	if err := doc.DataTo(&message); err != nil {
		return nil, err
	}

	// Set the ID from the document ID
	message.SetIDFromString(doc.Ref.ID)
	return &message, nil
}



func (r *firestoreMessageRepository) List(ctx context.Context, filter map[string]interface{}, sortField, sortOrder string, skip, limit int) ([]*models.Message, int, error) {
	query := r.client.Collection(r.collection).Query

	// Apply filters
	for key, value := range filter {
		if key == "_id" || key == "id" {
			// For Firestore, we need to use the document ID filter differently
			// This would require a separate query or handling
			continue
		}
		query = query.Where(key, "==", value)
	}

	// Sorting - default to timestamp descending for recent messages first
	direction := firestore.Desc
	if sortField == "" {
		sortField = "timestamp"
	}
	if sortOrder == "ASC" {
		direction = firestore.Asc
	}
	query = query.OrderBy(sortField, direction)

	// Get total count (this is a simplified approach, in production you might want to cache this)
	totalIter := query.Documents(ctx)
	totalCount := 0
	for {
		_, err := totalIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		totalCount++
	}

	// Apply pagination
	query = query.Offset(skip).Limit(limit)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var messages []*models.Message
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating documents: %v", err)
			return nil, 0, err
		}

		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			log.Printf("Error decoding document: %v", err)
			return nil, 0, err
		}

		// Set the ID from the document ID
		message.SetIDFromString(doc.Ref.ID)
		messages = append(messages, &message)
	}

	return messages, totalCount, nil
}

func (r *firestoreMessageRepository) FindByTopic(ctx context.Context, topic string, limit int) ([]*models.Message, error) {
	query := r.client.Collection(r.collection).
		Where("topic", "==", topic).
		OrderBy("timestamp", firestore.Desc).
		Limit(limit)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var messages []*models.Message
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, err
		}

		// Set the ID from the document ID
		message.SetIDFromString(doc.Ref.ID)
		messages = append(messages, &message)
	}

	return messages, nil
}

func (r *firestoreMessageRepository) FindByClientID(ctx context.Context, clientID string, limit int) ([]*models.Message, error) {
	query := r.client.Collection(r.collection).
		Where("client_id", "==", clientID).
		OrderBy("timestamp", firestore.Desc).
		Limit(limit)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var messages []*models.Message
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, err
		}

		// Set the ID from the document ID
		message.SetIDFromString(doc.Ref.ID)
		messages = append(messages, &message)
	}

	return messages, nil
}

func (r *firestoreMessageRepository) FindByTimeRange(ctx context.Context, from, to time.Time, limit int) ([]*models.Message, error) {
	query := r.client.Collection(r.collection).
		Where("timestamp", ">=", from).
		Where("timestamp", "<=", to).
		OrderBy("timestamp", firestore.Desc).
		Limit(limit)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var messages []*models.Message
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		var message models.Message
		if err := doc.DataTo(&message); err != nil {
			return nil, err
		}

		// Set the ID from the document ID
		message.SetIDFromString(doc.Ref.ID)
		messages = append(messages, &message)
	}

	return messages, nil
}
