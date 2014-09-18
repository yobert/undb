package undb

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
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

func (store* Store) ws_exec(ops []Op, source string) error {
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
	opchan := make(chan Op, 0)

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

	go func() {
		defer close(opchan)
		defer store.ws_cleanup(opchan) // double cleanup is ok

		for {
			_, msgdata, err := ws.ReadMessage()
			if err != nil {
				//fmt.Println(err.Error()) // don't really care
				return
			}

			var opslice []Op
			err = json.Unmarshal(msgdata, &opslice)
			if err != nil {
				panic(err)
			}

			err = store.ws_exec(opslice, ws_id)
			if err != nil {
				panic(err)
			}
		}
	}();

	for {
		op, ok := <-opchan
		if !ok {
			return
		}

		if op.changesource == ws_id {
			continue
		}

		opslice := []Op{op}
		j, err := json.Marshal(&opslice)
		if err != nil {
			panic(err)
		}
		err = ws.WriteMessage(websocket.TextMessage, j)
		if err != nil {
			panic(err)
		}
	}
}

