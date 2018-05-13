package datastore

import (
	level "github.com/syndtr/goleveldb/leveldb"
)

// Level provides datastore.Datastore interface over boltdb
type Level struct {
	*level.DB
}

func makeKey(r Record) []byte {
	return append([]byte(r.Name()), r.ID()...)
}

// Put stores Record in db
func (l *Level) Put(r Record) error {
	v, err := r.Marshal()
	if err != nil {
		return err
	}
	return l.DB.Put(makeKey(r), v, nil)
}

// Get retrieves record from db
func (l *Level) Get(r Record) error {
	buf, err := l.DB.Get(makeKey(r), nil)
	if err == nil {
		err = r.Unmarshal(buf)
	}
	return err
}

// Delete removes record form db
func (l *Level) Delete(r Record) error {
	return l.DB.Delete(makeKey(r), nil)
}
