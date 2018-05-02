// Package datastore provides interface and connectivity to message broker services
package datastore

// Data represents a struct which may be r/w to/from datastore
type Data interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Name() string // returns the structs name
}

// Query performs a read-only operation
type Query struct {
	Model string
	Args  map[string]interface{}
}

// Datastore provides generic interface for an API datastore
type Datastore interface {
	Create(*Query, Data) error
	Get(*Query) (Data, error)
	Update(*Query, Data) error
	Delete(*Query) error
}

// MockDB a simple API datastore implementation
type MockDB struct {
	DB map[string]map[string][]byte
}
