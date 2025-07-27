package undb

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func (store *Store) Save(path string) {
	log.Println("db save('" + path + "')")
	store.Lock()
	defer store.Unlock()

	var ops []Op
	store.Instruct(&ops)

	jsonbytes, err := json.Marshal(&ops)
	if err != nil {
		log.Println("save failed:", err)
		return
	}

	err = ioutil.WriteFile(path, jsonbytes, 0644)
	if err != nil {
		log.Println("save failed:", err)
		return
	}

	log.Println("db save('"+path+"') success", len(ops), "ops", len(jsonbytes), "bytes")
	store.dirty = false
}
