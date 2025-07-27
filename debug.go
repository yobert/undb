package undb

import (
	"bytes"
	"encoding/json"
)

func (s *Store) DebugString() string {
	b, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	return out.String()
}
