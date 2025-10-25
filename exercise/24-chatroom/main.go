package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var rooms = make(map[string]map[*websocket.Conn]string)
var mu sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

func wsChatRoom(c echo.Context) error {
	roomID := c.Param("roomId")
	username := c.Param("username")

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer conn.Close()

	mu.Lock()
	if rooms[roomID] == nil {
		rooms[roomID] = make(map[*websocket.Conn]string)
	}
	rooms[roomID][conn] = username
	mu.Unlock()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(rooms[roomID], conn)
			mu.Unlock()
			break
		}

		broadcastMsg := fmt.Sprintf("%s: %s", username, string(msg))

		mu.Lock()
		for cConn := range rooms[roomID] {
			if cConn != conn { 
				_ = cConn.WriteMessage(websocket.TextMessage, []byte(broadcastMsg))
			}
		}
		mu.Unlock()
	}

	return nil
}

func main() {
	e := echo.New()
	e.GET("/ws/chat/:roomId/user/:username", wsChatRoom)
	e.Logger.Fatal(e.Start(":8080"))
}
