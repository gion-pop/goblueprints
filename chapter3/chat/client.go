package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// client はチャットに参加してる 1 ユーザーを表す
type client struct {
	// socket はこの client のための WebSocket
	socket *websocket.Conn
	// send はメッセージが送られるチャネル
	send chan *message
	// room は client が現在参加してる部屋
	room *room
	// userData はユーザーに関する情報を保持
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
