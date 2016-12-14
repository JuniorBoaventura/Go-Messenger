package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "net/http"
  "crypto/md5"
  "encoding/hex"
  "sync"
)

var mu = &sync.Mutex{}
var clients = make(map[string]Client)

type Message struct {
    Name string
    Body string
    Type string
}

type User struct {
  Id string
  Name string
}

type Client struct {
	Websocket *websocket.Conn
  User User
}

type Users struct {
  Type string
  Users []User
}

// Send an object to all the clients connected except the current user
func (this User) broadcast(message interface{})  {
  mu.Lock()
  for _, client := range clients {
    if client.User.Id != this.Id {
      websocket.JSON.Send(client.Websocket, message)
    }
  }
  mu.Unlock()
}

// Send an object to all the clients connected
func (this Client) emit(message interface{})  {
  mu.Lock()
  for _, client := range clients {
    websocket.JSON.Send(client.Websocket, message)
  }
  mu.Unlock()
}

// Broadcast the deconnection of the user and remove it from clients
func (this User) disconnect() {
  disconnected := Message {this.Name, this.Id, "disconnected"}
  this.broadcast(disconnected)
  mu.Lock()
  delete(clients, this.Id)
  mu.Unlock()
}

// Return the list of the users connected
func (this Client) UsersConnected() []User {
  var connectedUsers []User
  mu.Lock()
  for _, client := range clients {
    if client.User.Id != this.User.Id {
      connectedUsers = append(connectedUsers, client.User)
    }
  }
  mu.Unlock()
  return connectedUsers
}

func main() {
  http.Handle("/ws", websocket.Handler(WebSocket))
  http.ListenAndServe(":8080", nil)
}

func WebSocket(ws *websocket.Conn) {

  ip := ws.Request().RemoteAddr
  clientId := hashMD5(ip)
  var client Client
  var user User
  fmt.Println("Client connected:", clientId)

  for {

    var data Message

    if err := websocket.JSON.Receive(ws, &data); err != nil {
      // If there is an error the user closed the connection
      fmt.Println("Client disconnected:", clientId)
      user.disconnect()
      return
    }

    if data.Type == "connect" {

      username := data.Body

      user = User{clientId, username}
      client = Client{ws, user}

      // User Receive is connection information
      connection := Message{username, clientId, "connected"}
      websocket.JSON.Send(ws, connection)

      // User receive the list of the connected users
      users := Users{"ConnectedUsers", client.UsersConnected()}
      websocket.JSON.Send(ws, users)

      // Prevent all the other clients of it's connection
      newUser := Message{username, clientId, "newUser"}
      user.broadcast(newUser)

      mu.Lock()
      clients[clientId] = client
      mu.Unlock()
    }

    if data.Type == "message" {
      client.emit(data)
    }

  }

}

func hashMD5(text string) string {
  hasher := md5.New()
  hasher.Write([]byte(text))
  return hex.EncodeToString(hasher.Sum(nil))
}
