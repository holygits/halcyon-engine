// Package datastore provides interface and connectivity to message broker services
package datastore

// Record represents a struct which may be r/w to/from datastore
type Record interface {
	// Marshal serialize to byte stream
	Marshal() ([]byte, error)
	// Deserialize from byte stream
	Unmarshal([]byte) error
	// ID returns the struct's unique ID
	ID() []byte
	// Name returns the struct type
	Name() string
}

// Datastore provides generic interface for an API datastore
type Datastore interface {
	Put(Record) error
	Get(Record) error
	Delete(Record) error
}
