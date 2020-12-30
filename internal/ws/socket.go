package ws

import (
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/domain"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
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
	AuthContext domain.WsAuth
	socket      *websocket.Conn
	findHandler FindHandler
}

// NewClient accepts a socket and returns an initialized Socket.
func NewWS(socket *websocket.Conn, auth domain.WsAuth, findHandler FindHandler) *Socket {
	return &Socket{
		socket:      socket,
		AuthContext: auth,
		findHandler: findHandler,
	}
}

// Write receives messages from the channel and writes to the socket.
func (c *Socket) Write(msg Message) {
	err := c.socket.WriteJSON(msg)
	if err != nil {
		c.logWithAuthContext().
			WithError(err).
			Warn("Socket write failed.")
	}
}

// Read intercepts messages on the socket and assigns them to a handler function.
func (c *Socket) Loop(onClose func()) {
	ppd := c.pingPong()
	defer func() {
		onClose()
		ppd <- true
		_ = c.socket.Close()
		c.logWithAuthContext().Info("Closing client")
	}()

	var msg Message
	for {
		if err := c.socket.ReadJSON(&msg); err != nil {
			closeErr, ok := err.(*websocket.CloseError)
			if ok && isGracefulClose(closeErr) {
				c.logWithAuthContext().
					Debug("Socket closed gracefully.")
			} else {
				c.logWithAuthContext().
					WithError(err).
					Warn("Socket read failed.")
			}
			break
		}

		if handler, found := c.findHandler(Event(msg.Name)); found {
			handler(c, msg.Data)
		} else {
			c.logWithAuthContext().
				WithField("message", msg).
				Warn("Failed to find handler for message.")
		}
	}
}

func isGracefulClose(err *websocket.CloseError) bool {
	return err.Code == websocket.CloseGoingAway || err.Code == websocket.CloseNormalClosure
}

func (c *Socket) logWithAuthContext() *logrus.Entry {
	return logrus.
		WithField("user", c.AuthContext.UserEmail).
		WithField("hdID", c.AuthContext.HdID)
}

func (c *Socket) pingPong() chan bool {
	done := make(chan bool)

	pongWait := viper.GetDuration(conf.WebSocketMaxPingPongDelay)
	err := c.socket.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		return nil
	}
	c.socket.SetPongHandler(func(string) error {
		c.logWithAuthContext().
			WithField("pongWait", pongWait).
			Debug("Pong received. Postponing timeout by pongWait.")
		return c.socket.SetReadDeadline(time.Now().Add(pongWait))
	})

	pingPeriod := viper.GetDuration(conf.WebSocketPingPongInterval)
	ticker := time.NewTicker(pingPeriod)

	go func() {
		c.logWithAuthContext().Debug("Starting ping pong.")
		for {
			select {
			case <-done:
				c.logWithAuthContext().Debug("Stopping ping pong.")
				ticker.Stop()
				return
			case <-ticker.C:
				err := c.socket.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
				if err != nil {
					c.logWithAuthContext().Warn("Failed to send ping.")
					return
				} else {
					c.logWithAuthContext().Debug("Sending ping.")
				}
			}
		}
	}()

	return done
}
