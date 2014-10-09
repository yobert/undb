package undb

import (
	"log"
	"bufio"
	"os"
	"encoding/json"
)

func (store *Store) Save(path string) {
	log.Println("db save('" + path + "')")
	store.Lock()
	defer store.Unlock()

	f, err := os.Create(path)
	if err != nil {
		log.Println("save failed: ", err)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	var ops []Op
	store.Instruct(&ops)

	jsonbytes, err := json.Marshal(&ops)
	if err != nil {
		log.Println("save failed: ", err)
		return
	}

	count, err := w.Write(jsonbytes)
	if err != nil {
		log.Println("save failed: ", err)
		return
	}

	err = w.Flush()
	if err != nil {
		log.Println("save failed: ", err)
		return
	}

	log.Println("db save('" + path + "') success", len(ops), "ops", count, "bytes")
}