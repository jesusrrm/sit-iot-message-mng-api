package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MessageType represents the type of IoT message based on MQTT topics
type MessageType string

const (
	MessageTypeStatus    MessageType = "status"    // Device status updates
	MessageTypeEvent     MessageType = "event"     // Device events/notifications
	MessageTypeOnline    MessageType = "online"    // Device online/offline status
	MessageTypeCommand   MessageType = "command"   // Device commands
	MessageTypeTelemetry MessageType = "telemetry" // Device telemetry data
	MessageTypeAlert     MessageType = "alert"     // Device alerts
	MessageTypeRPC       MessageType = "rpc"       // RPC calls
	MessageTypeUnknown   MessageType = "unknown"   // Unknown message type
)

// MessageStatus represents the processing status of a message
type MessageStatus string

const (
	MessageStatusReceived  MessageStatus = "received"  // Message received from MQTT
	MessageStatusProcessed MessageStatus = "processed" // Message successfully processed
	MessageStatusFailed    MessageStatus = "failed"    // Message processing failed
	MessageStatusPending   MessageStatus = "pending"   // Message waiting to be processed
)

// Message represents an IoT MQTT message in the system
type Message struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Topic      string                 `bson:"topic" json:"topic"`                     // MQTT topic
	Payload    string                 `bson:"payload" json:"payload"`                 // Raw message payload
	Timestamp  time.Time              `bson:"timestamp" json:"timestamp"`             // Message timestamp
	Marshalled map[string]interface{} `bson:"marshalled,omitempty" json:"marshalled"` // Parsed JSON payload (variable structure)
	ClientID   string                 `bson:"client_id" json:"clientId"`              // MQTT client ID (device identifier)

	// Derived/computed fields
	Type      MessageType   `bson:"type,omitempty" json:"type"`           // Message type derived from topic
	Status    MessageStatus `bson:"status,omitempty" json:"status"`       // Processing status
	DeviceID  string        `bson:"deviceId,omitempty" json:"deviceId"`   // Device ID (derived from client_id or topic)
	ProjectID string        `bson:"projectId,omitempty" json:"projectId"` // Associated project ID

	// Metadata and audit fields
	ProcessedAt *time.Time        `bson:"processedAt,omitempty" json:"processedAt"`     // When message was processed
	CreatedAt   time.Time         `bson:"createdAt" json:"createdAt"`                   // When record was created
	UpdatedAt   time.Time         `bson:"updatedAt" json:"updatedAt"`                   // When record was last updated
	CreatedBy   string            `bson:"createdBy,omitempty" json:"createdBy"`         // User who processed/created record
	Metadata    map[string]string `bson:"metadata,omitempty" json:"metadata,omitempty"` // Additional metadata
}

// Device represents an IoT device
type Device struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProjectID string             `bson:"projectId" json:"projectId"`
	Name      string             `bson:"name" json:"name"`
	Type      string             `bson:"type" json:"type"`
	Status    string             `bson:"status" json:"status"`
	ClientID  string             `bson:"clientId" json:"clientId"`           // MQTT client ID
	LastSeen  *time.Time         `bson:"lastSeen,omitempty" json:"lastSeen"` // Last time device sent a message
	Metadata  map[string]string  `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
	CreatedBy string             `bson:"createdBy" json:"createdBy"`
}

// MessageFilter represents filtering options for message queries
type MessageFilter struct {
	ProjectID    string        `json:"projectId,omitempty"`
	DeviceID     string        `json:"deviceId,omitempty"`
	ClientID     string        `json:"clientId,omitempty"`
	Type         MessageType   `json:"type,omitempty"`
	Status       MessageStatus `json:"status,omitempty"`
	TopicPattern string        `json:"topicPattern,omitempty"` // For filtering by topic pattern
	FromTime     *time.Time    `json:"fromTime,omitempty"`
	ToTime       *time.Time    `json:"toTime,omitempty"`
}

// GetMessageTypeFromTopic derives message type from MQTT topic
func GetMessageTypeFromTopic(topic string) MessageType {
	if topic == "" {
		return MessageTypeUnknown
	}

	// Check common patterns in the topic
	switch {
	case contains(topic, "/status/"):
		return MessageTypeStatus
	case contains(topic, "/events/"):
		return MessageTypeEvent
	case contains(topic, "/online"):
		return MessageTypeOnline
	case contains(topic, "/rpc"):
		return MessageTypeRPC
	case contains(topic, "/command"):
		return MessageTypeCommand
	case contains(topic, "/telemetry"):
		return MessageTypeTelemetry
	case contains(topic, "/alert"):
		return MessageTypeAlert
	default:
		return MessageTypeUnknown
	}
}

// GetDeviceIDFromClientID extracts device ID from client ID (can be customized based on naming convention)
func GetDeviceIDFromClientID(clientID string) string {
	// For now, use client ID as device ID
	// This can be customized based on your device naming convention
	return clientID
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				indexOf(s, substr) >= 0)))
}

// Helper function to find substring index
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
