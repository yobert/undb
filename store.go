package undb

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	//"log"
)

type StoreType int

func (s StoreType) String() string {
	switch s {
	case STORES:
		return "STORES"
	case VALUES:
		return "VALUES"
	}
	return "INVALID"
}

const (
	INVALIDSTORETYPE StoreType = iota
	STORES
	VALUES
)

type Store struct {
	Id      string
	Type    StoreType
	Records map[string]interface{}
	Deleted bool

	last   int
	parent *Store

	lock      sync.RWMutex `json:-`
	listeners map[chan Op]struct{}
	dirty     bool
}

func (store *Store) Lock() {
	store.lock.Lock()
}
func (store *Store) Unlock() {
	store.lock.Unlock()
}
func (store *Store) RLock() {
	store.lock.RLock()
}
func (store *Store) RUnlock() {
	store.lock.RUnlock()
}

func New(id string, typ StoreType) *Store {
	s := Store{
		Id:      id,
		Type:    typ,
		Records: make(map[string]interface{}),
	}
	return &s
}

func (store *Store) Insert(insert *Store, source string) error {
	_, exists := store.Records[insert.Id]
	if exists {
		return errors.New("Insert into store '" + store.Id + "' failed: '" + insert.Id + "' already exists")
	}

	// hack around incrementing .next when loading a db from disk.
	// ids generated in javascript will have strings in them so they're ok
	i, err := strconv.Atoi(insert.Id)
	if err == nil {
		if i > store.last {
			store.last = i
			//log.Println("insert on store " + store.Id + " increased last to", store.last)
		}
	}

	insert.parent = store
	store.Records[insert.Id] = insert
	store.Emit(&Op{Method: INSERT, Id: insert.Id, Type: insert.Type}, source)
	return nil
}

func (store *Store) Delete(source string) {
	store.Deleted = true
	store.Emit(&Op{Method: DELETE}, source)
}

func (store *Store) Update(values map[string]interface{}, source string) error {
	if store.Type != VALUES {
		return errors.New("Update store '" + store.Id + "' failed: not a VALUES store")
	}
	store.Records = values
	store.Emit(&Op{Method: UPDATE, Values: values}, source)
	return nil
}

func (store *Store) Merge(values map[string]interface{}, source string) error {
	if store.Type != VALUES {
		return errors.New("Merge store '" + store.Id + "' failed: not a VALUES store")
	}
	for k, v := range values {
		store.Records[k] = v
	}
	store.Emit(&Op{Method: MERGE, Values: values}, source)
	return nil
}

func (store *Store) Path() string {
	path := store.Id
	parent := store.parent
	for parent != nil {
		path = parent.Id + "." + path
		parent = parent.parent
	}
	return path
}

func (store *Store) FindOrPanic(path string) *Store {
	s := store.Find(path)
	if s == nil {
		panic("Find on store '" + store.Id + "' failed: Find path '" + path + "' failed")
	}
	return s
}

func (store *Store) Find(path string) *Store {
	chunks := strings.Split(path, ".")
	if len(chunks) < 1 || chunks[0] != store.Id {
		return nil
	}

	s := store
	for i, k := range chunks {
		if i == 0 {
			continue
		}

		if s.Type != STORES {
			return nil
		}
		v, ok := s.Records[k]
		if !ok {
			return nil
		}
		s, ok = v.(*Store)
		if !ok {
			return nil
		}
	}

	return s
}

func (store *Store) Seq() string {
	return strconv.Itoa(store.last + 1)
}
