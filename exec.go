package gorp

import (
	"errors"
)

type Op struct {
	Method string
	Name string
	Record interface{}
}

func (store *Store) Exec(op *Op) (interface{}, error) {
	switch op.Method {
	case "Insert":
		return nil, store.Insert(op.Name, op.Record)
	case "Update":
		return nil, store.Update(op.Name, op.Record)
	case "Upsert":
		store.Upsert(op.Name, op.Record)
		return nil, nil
	case "Delete":
		store.Delete()
		return nil, nil
	case "Get":
		v, exists := store.Get(op.Name)
		if !exists {
			return nil, nil
		}
		return v, nil
	case "Keys":
		v := store.Keys()
		return v, nil
	}

	return nil, errors.New("invalid method: '" + op.Method + "'")
}
