package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "net/http"
  "crypto/md5"
  "encoding/hex"
  "encoding/json"
)

var clients = make(map[string]Client)


type Message struct {
    Name string
    Body string
    Type string
}

type User struct {
  Id string  `json:"id"`
  Username string `json:"username"`
}

type Client struct {
	Websocket *websocket.Conn
  user User
}

type Users struct {
  Type string
  Users []User
}

func (this Client) broadcast(clients map[string]Client, message interface{})  {
  for _, client := range clients {
    if client.user.Id != this.user.Id {
      websocket.JSON.Send(client.Websocket, message)
    }
  }
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Echo(ws *websocket.Conn) {

  ip := ws.Request().RemoteAddr
  clientId := MD5Hash(ip)
  var client Client
  fmt.Println("Client connected:", clientId)

  for {

    var msg Message

    if err := websocket.JSON.Receive(ws, &msg); err != nil {
      fmt.Println(err)
      return
    }

    if msg.Type == "connect" {
      username := msg.Body

      user := User{clientId, username}
      client = Client{ws, user}



      connection := Message{username, clientId, "connected"}
      newUser := Message{username, clientId, "newUser"}

      var connectedUsers []User

      for _, client := range clients {
        connectedUsers = append(connectedUsers, client.user)
      }

      users := Users{"ConnectedUsers", connectedUsers}
      json, _ := json.Marshal(connectedUsers)
      fmt.Println(string(json))

      websocket.JSON.Send(ws, connection)
      websocket.JSON.Send(ws, users)

      client.broadcast(clients, newUser)
      clients[clientId] = client
    }

    if msg.Type == "message" {
      client.broadcast(clients, msg)
    }

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
