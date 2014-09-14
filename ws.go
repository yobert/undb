package undb

import (
	"github.com/gorilla/websocket"
	"encoding/json"
	"fmt"
	"io"
)

var next_ws_id = 0

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

func (store *Store) Websocket() func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		next_ws_id++;
		ws_id := fmt.Sprintf("ws-%d", next_ws_id);

		fmt.Println("Store '" + store.Name + "': websocket opened (" + ws_id + ")")
		defer ws.Close()
		defer fmt.Println("Store '" + store.Name + "': websocket closed")

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
					if err == websocket.ErrCloseSent {
						return
					}
					if err == io.EOF {
						return
					}
					if err == io.ErrUnexpectedEOF {
						return
					}
					panic(err)
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
}

/*func http_websocket_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Println("websocket opened")

	for {
		_, msgdata, err := ws.ReadMessage() // for now ignore whether its a binary or text message
		if err != nil {
			if err == io.EOF {
				log.Println("websocket closed")
				ws.Close()
				return
			}
			panic(err)
		}

		msg := WsMessage{}

		err = json.Unmarshal(msgdata, &msg)
		if err != nil {
			panic(err)
		}

		tmp, err := json.Marshal(&msg)
		ws.WriteMessage(websocket.TextMessage, tmp)
	}
}*/
