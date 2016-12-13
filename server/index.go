package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "net/http"
)

var clients []ClientConn

type Message struct {
    Name string
    Body string
}

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

func (this Message) GetName() string  {
  return this.Name
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Echo(ws *websocket.Conn) {

  ip := ws.Request().RemoteAddr
  fmt.Println("Client connected:", ip)
  client := ClientConn{ws, ip}

  clients = append(clients, client)

  for {
    var msg Message

    if err := websocket.JSON.Receive(ws, &msg); err != nil {
      fmt.Println(err)
      return
    }

    broadcastClients(msg)

  }

}

func broadcastClients(message Message) {
  for _, client := range clients {
    websocket.JSON.Send(client.websocket, message)
  }
}

func main() {
    http.HandleFunc("/", handler)
    http.Handle("/ws", websocket.Handler(Echo))
    http.ListenAndServe(":8080", nil)
}
