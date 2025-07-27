package undb

import (
	"log"
	"time"
)

func (store *Store) Snapshots(path string, sleep time.Duration) {
	log.Println("taking snapshots, interval", sleep)
	for {
		time.Sleep(sleep)
		if store.dirty {
			filename := path + time.Now().Local().Format("2006-01-02-15-04-05") + ".json"
			store.Save(filename)
		}
	}
}
