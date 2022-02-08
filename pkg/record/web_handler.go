package record

import (
	"log"
	"net/http"
)

// ServeWs handles websocket requests from the peer.
func (s *RecordService) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &WsConn{service: s, conn: conn, send: make(chan []byte, 256)}
	client.service.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}
