package undb

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	// "time"
)

func (store *Store) ws_setup(ops *[]Op, opchan chan Op) {
	store.Lock()
	defer store.Unlock()

	store.Instruct(ops)
	store.Listen(opchan)
}

func (store *Store) ws_cleanup(opchan chan Op) {
	store.Lock()
	defer store.Unlock()

	store.Unlisten(opchan)
}

func (store *Store) ws_exec(ops []Op, source string) error {
	store.Lock()
	defer store.Unlock()

	for _, op := range ops {
		err := store.Exec(&op, source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (store *Store) Websocket(ws *websocket.Conn, ws_id string) {
	fmt.Println("Store '" + store.Id + "': websocket opened (" + ws_id + ")")
	defer fmt.Println("Store '" + store.Id + "': websocket closing")

	var ops []Op
	opchan := make(chan Op, 8192)

	// this logic is tricky.  take care with these
	// locks and with defer.  lock/unlock is required for
	// store.Unlisten()
	store.ws_setup(&ops, opchan)
	defer store.ws_cleanup(opchan)

	j, err := json.Marshal(ops)
	if err != nil {
		panic(err)
	}
	ws.WriteMessage(websocket.TextMessage, j)

	//	time.Sleep(time.Second)

	go func() {
		defer close(opchan)
		defer store.ws_cleanup(opchan) // double cleanup is ok

		for {
			//			time.Sleep(time.Second)

			//log.Println("receiving message")
			_, msgdata, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//log.Println("receiving message complete")

			var opslice []Op
			err = json.Unmarshal(msgdata, &opslice)
			if err != nil {
				log.Println("message unmarshal error", err)
				continue
			}

			err = store.ws_exec(opslice, ws_id)
			if err != nil {
				panic(err)
			}
		}
	}()

	for {
		op, ok := <-opchan
		if !ok {
			return
		}

		if op.changesource == ws_id {
			continue
		}

		//		time.Sleep(time.Second)

		opslice := []Op{op}
		j, err := json.Marshal(&opslice)
		if err != nil {
			panic(err)
		}
		//log.Println("writing message")
		err = ws.WriteMessage(websocket.TextMessage, j)
		//log.Println("writing message complete")
		if err != nil {
			log.Println(err)
			return
		}
	}
}
