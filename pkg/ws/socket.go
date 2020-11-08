package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Message is a object used to pass data on sockets.
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

// Handler is a type representing functions which resolve requests.
type Handler func(*Socket, interface{})

func (m *Message) String() string {
	return fmt.Sprintf("message %s with data %+v", m.Name, m.Data)
}

// FindHandler is a type that defines handler finding functions.
type FindHandler func(Event) (Handler, bool)

// Socket is a type that reads and writes on sockets.
type Socket struct {
	socket      *websocket.Conn
	findHandler FindHandler
}

// NewClient accepts a socket and returns an initialized Socket.
func NewWS(socket *websocket.Conn, findHandler FindHandler) *Socket {
	return &Socket{
		socket:      socket,
		findHandler: findHandler,
	}
}

// Write receives messages from the channel and writes to the socket.
func (c *Socket) Write(msg Message) {
	err := c.socket.WriteJSON(msg)
	if err != nil {
		logrus.WithError(err).Warn("Socket write failed.")
	}
}

// Read intercepts messages on the socket and assigns them to a handler function.
func (c *Socket) InitReadLoop(onClose func()) {
	defer func() {
		onClose()
		_ = c.socket.Close()
		logrus.Info("Closing client")
	}()

	var msg Message
	for {
		if err := c.socket.ReadJSON(&msg); err != nil {
			logrus.WithError(err).Warn("Socket read failed.")
			break
		}

		if handler, found := c.findHandler(Event(msg.Name)); found {
			handler(c, msg.Data)
		} else {
			logrus.WithField("message", msg).Warn("Failed to find handler for message.")
		}
	}
}
