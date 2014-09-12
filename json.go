package gorp

import (
	"encoding/json"
)

type JsonError struct {
	Error string
}

func ErrorMarshal(e error) []byte {
	je := JsonError{e.Error()}
	b, err := json.Marshal(&je)
	if err != nil {
		panic(err)
	}
	return b
}

func (store *Store) ExecJson(input []byte) []byte {
	op := Op{}
	err := json.Unmarshal(input, &op)
	if err != nil {
		return ErrorMarshal(err)
	}

	r, e := store.Exec(&op)
	if e != nil {
		return ErrorMarshal(e)
	}

	b, be := json.Marshal(r)
	if be != nil {
		return ErrorMarshal(be)
	}
	return b
}
