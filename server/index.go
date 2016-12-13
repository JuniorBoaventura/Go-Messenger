package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "net/http"
  "crypto/md5"
  "encoding/hex"
)

var clients = make(map[string]ClientConn)

type Message struct {
    Name string
    Body string
    Type string
}

type ClientConn struct {
	websocket *websocket.Conn
	clientIP string
  Username string
}

func (this Message) GetName() string  {
  return this.Name
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Echo(ws *websocket.Conn) {

  ip := ws.Request().RemoteAddr
  clientId := MD5Hash(ip)

  fmt.Println("Client connected:", ip)

  for {

    var msg Message

    if err := websocket.JSON.Receive(ws, &msg); err != nil {
      fmt.Println(err)
      return
    }

    if msg.Type == "connect" {
      username := msg.Body

      client := ClientConn{ws, ip, username}
      clients[clientId] = client

      res := Message{username, clientId, "connected"}
      websocket.JSON.Send(ws, res)
    }

    if msg.Type == "message" {
      broadcastClients(msg)
    }

  }

}

func broadcastClients(message Message) {
  for _, client := range clients {
    websocket.JSON.Send(client.websocket, message)
  }
}

func MD5Hash(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func main() {
    http.HandleFunc("/", handler)
    http.Handle("/ws", websocket.Handler(Echo))
    http.ListenAndServe(":8080", nil)
}
