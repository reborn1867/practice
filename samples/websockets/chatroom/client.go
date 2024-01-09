package chatroom

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	message   chan []byte
	conn      *websocket.Conn
	RoomEntry *Room
}

func NewClient(room *Room, w http.ResponseWriter, r *http.Request) (*Client, error) {
	upgrader := websocket.Upgrader{
		WriteBufferSize: 1024,
		ReadBufferSize:  1024,
	}

	conn, err := upgrader.Upgrade(w, r, http.Header{})
	if err != nil {
		return nil, err
	}

	return &Client{
		message:   make(chan []byte, 256),
		conn:      conn,
		RoomEntry: room,
	}, nil
}

func (c *Client) Read() {
	defer func() {
		c.conn.Close()
		c.RoomEntry.leave <- c
	}()

	c.conn.SetReadLimit(512)
	err := c.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	if err != nil {
		logrus.Errorf("failed to set read deadline: %s", err)
	}
	c.conn.SetPongHandler(func(appData string) error {
		return c.conn.SetReadDeadline(time.Now().Add(10 * time.Minute))
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			logrus.Errorf("failed to read message: %s", err)
			break
		}
		c.RoomEntry.broadcast <- msg
	}
}

func (c *Client) Write() {
	// set ticker for ping message
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-c.message:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				logrus.Errorf("failed to send close messge: %s", err)
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				logrus.Errorf("failed to write message: %s", err)
			}
		case <-ticker.C:
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				logrus.Errorf("failed to ping: %s", err)
				return
			}
		}
	}
}

func ServeWs(room *Room, w http.ResponseWriter, r *http.Request) {
	logrus.Info("ws connection created")
	client, err := NewClient(room, w, r)
	if err != nil {
		logrus.Errorf("failed to new client: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	room.join <- client

	go client.Read()
	go client.Write()
}
