package main

import (
	"log"
	"net/http"

	"github.com/gion-xy/goblueprints/chapter1/trace"
	"github.com/gorilla/websocket"
)

type room struct {
	// forward は他の client に転送するためのメッセージを保持するチャネル
	forward chan []byte
	// join は room に参加しようとしている client のためのチャネル
	join chan *client
	// leave は room から退室しようとしているクライアントのためのチャネル
	leave chan *client
	// clients には在室しているすべての client が保持される
	clients map[*client]bool
	// tracer はチャットルームで行われた操作のログを受け取る
	tracer trace.Tracer
}

func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.Off(),
	}
}

func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
			r.tracer.Trace("新しいクライアントが参加")
		case client := <-r.leave:
			// 退室
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("クライアントが退室")
		case msg := <-r.forward:
			// room にいるすべての client にメッセージを転送
			for client := range r.clients {
				select {
				case client.send <- msg:
					// 送信成功
					r.tracer.Trace(" -- クライアントに送信")
				default:
					// 送信失敗
					delete(r.clients, client)
					close(client.send)
					r.tracer.Trace(" -- 送信に失敗")
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 512
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()
}
