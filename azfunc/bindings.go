package azfunc

import (
	"net/http"
	"reflect"

	"github.com/Azure/azure-functions-go-worker/internal/logger"
	"github.com/Azure/azure-functions-go-worker/internal/rpc"
)

// TriggerType represents the supported trigger types.
type TriggerType string

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

var StringToType = map[string]reflect.Type{
	"*http.Request":          reflect.TypeOf((*http.Request)(nil)),
	"*azfunc.Context":        reflect.TypeOf((*Context)(nil)),
	"*azfunc.Blob":           reflect.TypeOf((*Blob)(nil)),
	"*azfunc.Timer":          reflect.TypeOf((*Timer)(nil)),
	"*azfunc.QueueMsg":       reflect.TypeOf((*QueueMsg)(nil)),
	"*azfunc.EventGridEvent": reflect.TypeOf((*EventGridEvent)(nil)),
	"map[string]interface{}": reflect.TypeOf(reflect.TypeOf((map[string]interface{})(nil))),
}

// Func contains a function symbol with in and out param types
type Func struct {
	Func             reflect.Value
	Bindings         map[string]*rpc.BindingInfo
	In               []reflect.Type
	NamedInArgs      []*Arg
	Out              []reflect.Type
	NamedOutBindings map[string]reflect.Value
}

// Context contains the runtime context of the function
type Context struct {
	FunctionID   string
	InvocationID string
	Logger       *logger.Logger
}

// Arg represents an initial representation of a func argument
type Arg struct {
	Name string
	Type reflect.Type
}

type Timer struct {
	PastDue bool `json:"IsPastDue"`
}

type QueueMsg struct {
	ID              string
	Data            string
	DequeueCount    int
	InsertionTime   string
	ExpirationTime  string
	NextVisibleTime string
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

// Blob contains the data from a blob as string
type Blob struct {
	Name   string
	URI    string
	Data   string
	Length int
}
