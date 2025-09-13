package models

import (
	"time"
)

// ClientAggregations holds all aggregation data for a client_id
type ClientAggregations struct {
	ClientID     string                                                      `bson:"client_id" json:"client_id"`
	Aggregations map[string]map[string]map[string]map[string]*AggregatedData `bson:"aggregations" json:"aggregations"`
	// Structure: channel -> variable -> period -> timestamp -> AggregatedData
}

// AggregatedData holds the aggregation result for a client/channel/variable/time period
type AggregatedData struct {
	ClientID  string    `bson:"client_id" json:"client_id"`
	Channel   string    `bson:"channel" json:"channel"`
	Variable  string    `bson:"variable" json:"variable"`
	Period    string    `bson:"period" json:"period"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Sum       float64   `bson:"sum" json:"sum"`
	Count     int       `bson:"count" json:"count"`
	Min       float64   `bson:"min" json:"min"`
	Max       float64   `bson:"max" json:"max"`
	Avg       float64   `bson:"avg" json:"avg"`
}