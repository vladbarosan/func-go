package azure

// Blob contains the data from a blob as string
type Blob struct {
	Name   string
	URI    string
	Data   string
	Length int
}
