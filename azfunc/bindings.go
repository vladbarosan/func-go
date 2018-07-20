package azfunc

import (
	"context"

	"github.com/Azure/azure-functions-go-worker/internal/logger"
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
	PastDue       bool           `json:"IsPastDue"`
	ScheduleStats ScheduleStatus `json:"ScheduleStatus"`
}

// ScheduleStatus contains the schedule for a Timer
type ScheduleStatus struct {
	Next        string `json:"Next"`
	Last        string `json:"Last"`
	LastUpdated string `json:"LastUpdated"`
}

// QueueMsg represents an Azure queue message
type QueueMsg struct {
	Text         string `json:"azfuncdata"`
	ID           string `json:"Id"`
	Insertion    string `json:"InsertionTime"`
	Expiration   string `json:"ExpirationTime"`
	PopReceipt   string `json:"PopReceipt"`
	NextVisible  string `json:"NextVisibleTime"`
	DequeueCount int    `json:"DequeueCount"`
}

// Blob contains the data from a blob as string
type Blob struct {
	Content    string         `json:"azfuncdata"`
	Name       string         `json:"name"`
	URI        string         `json:"Uri"`
	Properties BlobProperties `json:"Properties"`
}

// BlobProperties contain metadata about a blob
type BlobProperties struct {
	Length       int    `json:"Length"`
	ContentMD5   string `json:"ContentMD5"`
	ContentType  string `json:"ContentType"`
	ETag         string `json:"ETag"`
	LastModified string `json:"LastModified"`
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

// EventHubEvent properties of an event sent to an Event Hub.
type EventHubEvent struct {
	Data            string                 `json:"azfuncdata"`
	PartitionKey    *string                `json:"PartitionKey"`
	SequenceNumber  int                    `json:"SequenceNumber"`
	Offset          int                    `json:"Offset"`
	EnqueuedTimeUtc string                 `json:"EnqueuedTimeUtc"`
	Properties      map[string]interface{} `json:"Properties"`
}
