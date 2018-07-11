package azfunc

import (
	"context"

	"github.com/Azure/azure-functions-go-worker/internal/logger"
)

// BindingType represents the supported binding types.
type BindingType string

const (
	// HTTPTrigger represents a HTTP trigger in function load request from the host
	HTTPTrigger BindingType = "httpTrigger"

	// BlobTrigger represents a blob trigger in function load request from host
	BlobTrigger BindingType = "blobTrigger"

	// QueueTrigger represents a queue trigger in function load request from host
	QueueTrigger BindingType = "queueTrigger"

	// TimerTrigger represents a queue trigger in function load request from host
	TimerTrigger BindingType = "timerTrigger"

	// EventGridTrigger represents a queue trigger in function load request from host
	EventGridTrigger BindingType = "eventGridTrigger"

	// HTTPBinding represents a HTTP binding in function load request from the host
	HTTPBinding BindingType = "http"

	// BlobBinding represents a blob binding in function load request from the host
	BlobBinding BindingType = "blob"

	// QueueBinding represents a queue binding in function load request from the host
	QueueBinding BindingType = "queue"

	// TableBinding represents a table binding in function load request from the host
	TableBinding BindingType = "table"
)

// Context contains the runtime context of the function
type Context struct {
	context.Context
	FunctionID   string
	InvocationID string
	Logger       *logger.Logger
}

// Timer represents a timer trigger
type Timer struct {
	PastDue bool `json:"IsPastDue"`
}

// QueueMsg represents an Azure queue message
type QueueMsg struct {
	Text         string `json:"data"`
	ID           string `json:"Id"`
	Insertion    string `json:"InsertionTime"`
	Expiration   string `json:"ExpirationTime"`
	PopReceipt   string `json:"PopReceipt"`
	NextVisible  string `json:"NextVisibleTime"`
	DequeueCount int    `json:"DequeueCount"`
}

// Blob contains the data from a blob as string
type Blob struct {
	Name   string
	URI    string
	Data   string
	Length int
}

//EventGridEvent properties of an event published to an Event Grid topic.
type EventGridEvent struct {
	// ID - An unique identifier for the event.
	ID string `json:"id"`
	// Topic - The resource path of the event source.
	Topic string `json:"topic"`
	// Subject - A resource path relative to the topic path.
	Subject string `json:"subject"`
	// Data - Event data specific to the event type.
	Data map[string]interface{} `json:"data"`
	// EventType - The type of the event that occurred.
	EventType string `json:"eventType"`
	// EventTime - The time (in UTC) the event was generated.
	EventTime string `json:"eventTime"`
	// MetadataVersion - The schema version of the event metadata.
	MetadataVersion string `json:"metadataVersion"`
	// DataVersion - The schema version of the data object.
	DataVersion string `json:"dataVersion"`
}
