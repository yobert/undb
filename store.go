package gorp

import (
	"errors"
)

type Store struct {
	Name string
	Records map[string]interface{}
	Deleted bool
}

// Init() initializes a new root database object
func Init(name string) *Store {
	s := Store{
		Name: name,
		Records: make(map[string]interface{}),
	}
	return &s
}

func (store *Store) Insert(name string, record interface{}) error {
	_, exists := store.Records[name]
	if exists {
		return errors.New("Insert into store '" + store.Name + "' failed: '" + name + "' already exists")
	}
	store.Records[name] = record
	return nil
}

func (store *Store) Update(name string, record interface{}) error {
	_, exists := store.Records[name]
	if !exists {
		return errors.New("Update store '" + store.Name + "' failed: '" + name + "' does not exist")
	}
	store.Records[name] = record
	return nil
}

func (store *Store) Upsert(name string, record interface{}) {
	store.Records[name] = record
}

func (store *Store) Delete() {
	store.Deleted = true
}

func (store *Store) Get(name string) (interface{}, bool) {
	v, exists := store.Records[name]
	return v, exists
}

func (store *Store) Keys() []string {
	v := make([]string, len(store.Records))
	i := 0
	for k, _ := range store.Records {
		v[i] = k
	}
	return v
}
