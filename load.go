package undb

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

func Load(path string) (*Store, error) {
	log.Println("db load('" + path + "')")
	jsonbytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("load failed: ", err)
		return nil, err
	}

	source := "load"

	db := New("db", STORES)

	var ops []Op
	err = json.Unmarshal(jsonbytes, &ops)
	if err != nil {
		log.Println("load failed: ", err)
		return nil, err
	}

	for _, op := range ops {
		err := db.Exec(&op, source)
		if err != nil {
			log.Println("load failed: ", err)
			return nil, err
		}
	}

	log.Println("db load('" + path + "') success", len(ops), "ops", len(jsonbytes), "bytes")
	return db, nil
}