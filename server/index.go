package main

import (
  "golang.org/x/net/websocket"
  "fmt"
  "net/http"
  "encoding/json"
)

var clients []websocket.Conn

type Message struct {
    Name string
    Body string
}

func (this Message) GetName() string  {
  return this.Name
}

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func Echo(ws *websocket.Conn) {
  
  for {
    var msg []byte

    if err := websocket.Message.Receive(ws, &msg); err != nil {
      return
    }

    m := Message{"Alice", "Hello"}
    b, err := json.Marshal(m)

    if (err == nil) {
      fmt.Printf("hello")
      websocket.Message.Send(ws, string(b))
    }

    fmt.Printf(string(msg))
  }

}

func main() {
    http.HandleFunc("/", handler)
    http.Handle("/ws", websocket.Handler(Echo))
    http.ListenAndServe(":8080", nil)
}
