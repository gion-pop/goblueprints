package main

import (
	"github.com/gorilla/websocket"
)

// client はチャットに参加してる 1 ユーザーを表す
type client struct {
	// socket はこの client のための WebSocket
	socket *websocket.Conn
	// send はメッセージが送られるチャネル
	send chan []byte
	// room は client が現在参加してる部屋
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
