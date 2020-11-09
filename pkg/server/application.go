package server

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/api"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/service"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/ws"
	"github.com/gorilla/mux"
)

func RegisterREST(r *mux.Router, s *service.SlackService) {
	handlers := api.NewRestHandlers(s)

	r.HandleFunc("/rest/v1/happiness-door", handlers.Initiation).
		Methods("POST")
	r.HandleFunc("/rest/v1/happiness-door/interaction", handlers.Vote).
		Methods("POST")
}

func RegisterWS(r *mux.Router, s *service.SlackService) {
	wsRouter := ws.NewRouter(s)
	handlers := api.NewWSHandlers(s)
	wsRouter.Handle(handlers.CreateVoteHandler())
	r.Handle("/ws/v1/connect", wsRouter)
}
