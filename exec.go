package undb

import (
	"encoding/json"
	"errors"
)

type OpMethodType int

func (method OpMethodType) String() string {
	switch method {
	case INSERT:
		return "INSERT"
	case DELETE:
		return "DELETE"
	case UPDATE:
		return "UPDATE"
	case MERGE:
		return "MERGE"
	}
	return "INVALID"
}

const (
	INVALIDOPMETHOD OpMethodType = iota
	INSERT
	DELETE
	UPDATE
	MERGE
)

type Op struct {
	Method OpMethodType
	Path   string

	Id     string                 `json:"Id,omitempty"`
	Type   StoreType              `json:"Type,omitempty"`
	Values map[string]interface{} `json:"Values,omitempty"`

	changesource string
}

func (op *Op) Copy() (out Op) {
	// TODO a more efficient deep copy
	j, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(j, &out)
	if err != nil {
		panic(err)
	}
	return
}

func (store *Store) Exec(op *Op, source string) error {
	s := store.Find(op.Path)
	if s == nil {
		return errors.New("Exec on store '" + store.Id + "' failed: Find path '" + op.Path + "' failed")
	}

	switch op.Method {
	case INSERT:
		return s.Insert(New(op.Id, op.Type), source)
	case DELETE:
		s.Delete(source)
		return nil
	case UPDATE:
		return s.Update(op.Values, source)
	case MERGE:
		return s.Merge(op.Values, source)
	}
	return errors.New("Exec on store '" + store.Id + "' failed: Invalid method: '" + op.Method.String() + "'")
}
