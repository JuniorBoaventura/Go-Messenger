package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "net/http"
  "crypto/md5"
  "encoding/hex"
)

var clients = make(map[string]Client)

type Message struct {
    Name string
    Body string
    Type string
}

type User struct {
  Id string  `json:"id"`
  Name string `json:"username"`
}

type Client struct {
	Websocket *websocket.Conn
  user User
}

type Users struct {
  Type string
  Users []User
}

func (this User) broadcast(message interface{})  {
  for _, client := range clients {
    if client.user.Id != this.Id {
      websocket.JSON.Send(client.Websocket, message)
    }
  }
}

func (this Client) emit(message interface{})  {
  for _, client := range clients {
    websocket.JSON.Send(client.Websocket, message)
  }
}

func (this User) disconnect() {
  disconnected := Message {this.Name, this.Id, "disconnected"}
  this.broadcast(disconnected)
  delete(clients, this.Id)
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (this Client) UsersConnected() []User {
  var connectedUsers []User

  for _, client := range clients {
    if client.user.Id != this.user.Id {
      connectedUsers = append(connectedUsers, client.user)
    }
  }

  return connectedUsers
}

func Echo(ws *websocket.Conn) {

  ip := ws.Request().RemoteAddr
  clientId := MD5Hash(ip)
  var client Client
  var user User
  fmt.Println("Client connected:", clientId)

  for {

    var msg Message

    if err := websocket.JSON.Receive(ws, &msg); err != nil {
      user.disconnect()
      return
    }

    if msg.Type == "connect" {
      username := msg.Body

      user = User{clientId, username}
      client = Client{ws, user}

      connection := Message{username, clientId, "connected"}
      newUser := Message{username, clientId, "newUser"}

      users := Users{"ConnectedUsers", client.UsersConnected()}

      // User Receive is connection information
      websocket.JSON.Send(ws, connection)

      // User receive the list of the connected users
      websocket.JSON.Send(ws, users)

      // Prevent all the other clients of it's connection
      user.broadcast(newUser)

      clients[clientId] = client
    }

    if msg.Type == "message" {
      client.emit(msg)
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
