package chatroom

import (
	"github.com/sirupsen/logrus"
)

type Room struct {
	clients   map[*Client]struct{}
	join      chan *Client
	leave     chan *Client
	broadcast chan []byte
	quit      chan struct{}
}

func NewRoom() *Room {
	return &Room{
		clients:   make(map[*Client]struct{}),
		join:      make(chan *Client),
		leave:     make(chan *Client),
		broadcast: make(chan []byte),
		quit:      make(chan struct{}),
	}
}

func (r *Room) Start() {
	for {
		select {
		case c := <-r.join:
			r.clients[c] = struct{}{}
			logrus.Infof("There are %d users connecting", len(r.clients))
		case c := <-r.leave:
			delete(r.clients, c)
			close(c.message)
		case msg := <-r.broadcast:
			for c := range r.clients {
				c.message <- msg
			}
		case <-r.quit:
			logrus.Infof("closing chatroom")
			return
		}
	}
}
