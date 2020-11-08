package ws

import (
	"errors"
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/service"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const happinessDoorData = "HappinessDoorData"

// Event is a type representing request names.
type Event string

// Client represents a connected WS client
type Client struct {
	socket *Socket
	hdID   int
}

// Router is a message routing object mapping events to function handlers.
type Router struct {
	rules   map[Event]Handler
	service *service.SlackService
}

// NewRouter returns an initialized Router.
func NewRouter(s *service.SlackService) *Router {
	return &Router{
		rules:   make(map[Event]Handler),
		service: s,
	}
}

// ServeHTTP creates the socket connection and begins the read routine.
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hdId, proceed := rt.authRequest(w, r)
	if !proceed {
		return
	}

	rt.establishWs(hdId, w, r)
}

func (rt *Router) authRequest(w http.ResponseWriter, r *http.Request) (int, bool) {
	token, err := extractQueryParameter("t", r)
	if err != nil {
		logrus.WithError(err).Warn("Failed to authenticate request.")
		w.WriteHeader(401)
		return 0, false
	}
	hdId, err := extractQueryParameter("i", r)
	if err != nil {
		logrus.WithError(err).
			WithField("token", token).
			Warn("Failed to authenticate request.")
		w.WriteHeader(401)
		return 0, false
	}
	verified, err := rt.service.VerifyToken(hdId, token)
	if err != nil {
		logrus.WithError(err).
			WithField("token", token).
			WithField("hdId", hdId).
			Error("Token verification failed.")
		w.WriteHeader(500)
		return 0, false
	}
	if !verified {
		logrus.WithError(err).
			WithField("token", token).
			WithField("hdId", hdId).
			Warn("Provided token is not valid.")
		w.WriteHeader(401)
		return 0, false
	}
	id, _ := strconv.Atoi(hdId)
	return id, true
}

func (rt *Router) establishWs(hdID int, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// accept all?
		CheckOrigin: func(r *http.Request) bool {
			//TODO add check, check url against configuration
			return true
		},
	}

	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.WithError(err).Error("Socket server configuration error.")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ws := NewWS(socket, rt.findHandler)
	err = rt.sendInitData(ws, hdID)
	if err != nil {
		logrus.WithError(err).Warn("Failed to send initial data to socket.")
	}

	feed := rt.service.SubscribeHappinessDoorFeed(hdID)
	go func() {
		for hd := range feed {
			ws.Write(Message{Name: happinessDoorData, Data: hd})
		}
	}()
	ws.InitReadLoop(func() {
		rt.service.UnsubscribeHappinessDoorFeed(hdID, feed)
	})
}

func (rt *Router) sendInitData(s *Socket, hdID int) error {
	hd, err := rt.service.ComputeVoting(hdID)
	if err != nil {
		return err
	}
	s.Write(Message{Name: happinessDoorData, Data: hd})
	return nil
}

func (rt *Router) findHandler(event Event) (Handler, bool) {
	handler, found := rt.rules[event]
	return handler, found
}

// Handle is a function to add handlers to the router.
func (rt *Router) Handle(event Event, handler Handler) {
	// store in to router rules
	rt.rules[event] = handler
}

func extractQueryParameter(name string, r *http.Request) (string, error) {
	ts := r.URL.Query()[name]
	if len(ts) == 0 {
		return "", errors.New(fmt.Sprintf("param %v not present in request", name))
	}
	if len(ts) > 1 {
		return "", errors.New(fmt.Sprintf("multiple %v params present in request", name))
	}
	return ts[0], nil
}
