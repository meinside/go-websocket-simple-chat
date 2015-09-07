/*
	Simple WebSocket Server with HTML client

	created on : 2015.09.04.
	last update: 2015.09.07.

	meinside@gmail.com
*/

package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

const (
	Port        = 8080
	BufferSize  = 1024
	NumChannels = 5
)

/*
	{
		"user": "some_user",
		"message": "message from some_user"
	}
*/
type Message struct {
	User    string `json:user`
	Message string `json:message`
}

// connection channels
var cConn = make(chan *websocket.Conn, NumChannels)
var cDisconn = make(chan *websocket.Conn, NumChannels)

// data channels
var cData = make(chan []byte, NumChannels)

/*
	WebSocket handler (broadcasts received json data back to connected clients)
*/
func wsHandler(ws *websocket.Conn) {
	cConn <- ws

	read := make([]byte, BufferSize)
	for {
		if n, err := ws.Read(read); err == nil {
			//fmt.Printf("> Read %d byte(s): %s\n", n, string(read[:n]))
			cData <- read[:n]
		} else {
			//fmt.Printf("* Error: %s\n", err.Error())
			break
		}
	}

	cDisconn <- ws

	ws.Close()
}

func main() {
	go func() {
		var m Message
		conns := make(map[*websocket.Conn]struct{}) // empty connections' set
		for {
			select {
			case conn := <-cConn:
				conns[conn] = struct{}{} // add connection to the set

				fmt.Printf("> connected (total %d connections)\n", len(conns))
			case conn := <-cDisconn:
				if _, ok := conns[conn]; ok { // check if it exists in the set
					delete(conns, conn) // remove connection from the set

					fmt.Printf("> diconnected (total %d connections)\n", len(conns))
				}
			case data := <-cData:
				if err := json.Unmarshal(data, &m); err == nil { // when it is in JSON format,
					fmt.Printf("> [%s] %s\n", m.User, m.Message)

					// broadcast to all connections
					for conn, _ := range conns {
						conn.Write(data)
					}
				} else { // when it is not in JSON format,
					fmt.Printf("> %s\n", string(data))
				}
			}
		}
	}()

	http.Handle("/", http.FileServer(http.Dir("./web"))) // for serving HTML
	http.Handle("/server", websocket.Handler(wsHandler)) // for WebSocket

	fmt.Printf("> Listening on port: %d\n", Port)
	fmt.Printf("> Open 'http://localhost:%d/' in your web browser.\n", Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", Port), nil); err != nil {
		panic("Error: " + err.Error())
	}
}
