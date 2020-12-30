package ws

import (
	"errors"
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/domain"
	"github.com/Skamaniak/happiness-door-slack-bot/internal/service"
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
	auth, proceed := rt.authRequest(w, r)
	if !proceed {
		return
	}
	rt.establishWs(auth, w, r)
}

func (rt *Router) authRequest(w http.ResponseWriter, r *http.Request) (domain.WsAuth, bool) {
	auth, err := rt.extractAuth(r)
	if err != nil {
		logrus.WithError(err).
			WithField("auth", auth).
			Warn("Failed to authenticate request.")
		w.WriteHeader(401)
		return auth, false
	}

	verified, err := rt.verifyAuth(auth)
	if err != nil {
		logrus.WithError(err).
			WithField("auth", auth).
			Error("Auth verification failed.")
		w.WriteHeader(500)
		return auth, false
	}
	if !verified {
		logrus.WithError(err).
			WithField("auth", auth).
			Warn("Provided auth details are not valid.")
		w.WriteHeader(401)
		return auth, false
	}

	return auth, true
}

func (rt *Router) extractAuth(r *http.Request) (domain.WsAuth, error) {
	res := domain.WsAuth{}
	t, err := extractQueryParameter("t", r)
	if err != nil {
		return res, err
	}
	res.Token = t

	id, err := extractQueryParameter("i", r)
	if err != nil {
		return res, err
	}
	res.HdID, err = strconv.Atoi(id)
	if err != nil {
		return res, err
	}

	ue, err := extractQueryParameter("u", r) //TODO this will be Google OAuth JWT in the future
	if err != nil {
		return res, err
	}
	res.UserEmail = ue

	return res, nil
}

func (rt *Router) verifyAuth(auth domain.WsAuth) (bool, error) {
	verified, err := rt.service.VerifyToken(auth.HdID, auth.Token)
	if err != nil || !verified {
		return verified, err
	}
	return rt.service.VerifySlackUser(auth.UserEmail), nil
}

func (rt *Router) establishWs(auth domain.WsAuth, w http.ResponseWriter, r *http.Request) {
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

	ws := NewWS(socket, auth, rt.findHandler)
	err = rt.sendInitData(ws)
	if err != nil {
		logrus.WithError(err).Warn("Failed to send initial data to socket.")
	}

	feed := rt.service.SubscribeHappinessDoorFeed(auth.HdID)
	go func() {
		for hd := range feed {
			ws.Write(Message{Name: happinessDoorData, Data: hd})
		}
	}()
	ws.Loop(func() {
		rt.service.UnsubscribeHappinessDoorFeed(auth.HdID, feed)
	})
}

func (rt *Router) sendInitData(s *Socket) error {
	hd, err := rt.service.ComputeVoting(s.AuthContext.HdID)
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
