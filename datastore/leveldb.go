package datastore

import (
	level "github.com/syndtr/goleveldb/leveldb"
)

// Level provides datastore.Datastore interface over boltdb
type Level struct {
	*level.Storage
}

func makeKey(r Record) []byte {
	return append([]byte(r.Name()), r.ID()...)
}
func (db *level) Put(r Record) error {
	v, err := r.Marshal()
	if err != nil {
		return err
	}
	return db.Put(makeKey(r), v, nil)
}

func (db *Level) Get(r Record) error {
	buf, err := db.Get(makeKey(r), nil)
	if err == nil {
		err = r.Unmarshal(buf)
	}
	return err
}

func (db *Level) Delete(r Record) error {
	return db.Delete(makeKey(r), nil)
}
