package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageType represents the type of IoT message
type MessageType string

const (
	MessageTypeTelemetry MessageType = "telemetry"
	MessageTypeCommand   MessageType = "command"
	MessageTypeEvent     MessageType = "event"
	MessageTypeAlert     MessageType = "alert"
)

// MessageStatus represents the status of a message
type MessageStatus string

const (
	MessageStatusPending   MessageStatus = "pending"
	MessageStatusDelivered MessageStatus = "delivered"
	MessageStatusFailed    MessageStatus = "failed"
	MessageStatusProcessed MessageStatus = "processed"
)

// Message represents an IoT message in the system
type Message struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID string             `bson:"projectId" json:"projectId"`
	DeviceID  string             `bson:"deviceId" json:"deviceId"`
	Type      MessageType        `bson:"type" json:"type"`
	Status    MessageStatus      `bson:"status" json:"status"`
	Payload   interface{}        `bson:"payload" json:"payload"`
	Metadata  map[string]string  `bson:"metadata,omitempty" json:"metadata,omitempty"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
}

// Device represents an IoT device
type Device struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID string             `bson:"projectId" json:"projectId"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Status    string             `bson:"status" json:"status"`
	LastSeen  *time.Time         `bson:"lastSeen,omitempty" json:"lastSeen,omitempty"`
	Metadata  map[string]string  `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
}
