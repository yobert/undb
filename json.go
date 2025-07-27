package undb

import (
	"encoding/json"
	"errors"
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

/*func (store *Store) ExecJson(input []byte) []byte {
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
}*/

// build a sequence of ops to re-create the store.
// pass a *[]Op for it to append into
func (store *Store) Instruct(buf *[]Op) error {
	if buf == nil {
		return errors.New("Instruct() called with nil Op buffer")
	}

	path := store.Path()

	if store.Type == STORES {
		for k, v := range store.Records {
			vs, ok := v.(*Store)
			if !ok {
				return errors.New("Record '" + k + "' in store '" + path + "' is not a Store{}")
			}
			o := Op{
				Method: INSERT,
				Path:   path,
				Id:     vs.Id,
				Type:   vs.Type,
			}
			*buf = append(*buf, o)

			e := vs.Instruct(buf)
			if e != nil {
				return e
			}
		}
	} else if store.Type == VALUES {
		o := Op{
			Method: UPDATE,
			Path:   path,
			Values: store.Records,
		}
		*buf = append(*buf, o)
	} else {
		return errors.New("Invalid store type: " + store.Type.String())
	}

	if store.Deleted {
		o := Op{
			Method: DELETE,
			Path:   path,
		}
		*buf = append(*buf, o)
	}

	return nil
}
