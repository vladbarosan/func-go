package azure

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
