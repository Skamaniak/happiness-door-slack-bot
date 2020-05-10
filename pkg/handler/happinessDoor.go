package handler

import (
	"bytes"
	"encoding/json"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/db"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"github.com/slack-go/slack"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type action struct {
	Identifier string `json:"value"`
	Action     string `json:"action_id"`
}

type Handlers struct {
	repo *db.HappinessDoor
}

type interactiveResponse struct {
	ResponseUrl string   `json:"response_url"`
	Actions     []action `json:"actions"`
}

func NewHandlers(repo *db.HappinessDoor) *Handlers {
	return &Handlers{repo: repo}
}

func toJson(v interface{}) []byte {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		log.Println("WARN: failed to marshal slack message to JSON")
	}
	return jsonBytes
}

func logRequest(r *http.Request) {
	if requestBytes, err := httputil.DumpRequest(r, true); err != nil {
		log.Println("Failed to parse request", err)
	} else {
		log.Println(string(requestBytes))
	}
}

func writeResponse(response slack.Msg, w http.ResponseWriter) error {
	w.Header().
		Set("Content-Type", "application/json")

	jsonBytes := toJson(response)
	_, err := w.Write(jsonBytes)
	return err
}

func (h *Handlers) Initiation(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	defer func() { _ = r.Body.Close() }()

	slash, err := slack.SlashCommandParse(r)
	if err != nil {
		log.Println("WARN: Failed to parse request", err)
		return
	}

	meetingName := slash.Text
	id, err := h.repo.CreateHappinessDoor(meetingName)
	if err != nil {
		log.Println("WARN: Failed to create new happiness door record in db", err)
		return
	}

	message := domain.CreateSlackMessage(domain.StubRecord(id, meetingName))
	err = writeResponse(message, w)
	if err != nil {
		log.Println("WARN: Failed to respond to request", err)
	}

}

func (h *Handlers) incrementVoting(result interactiveResponse) (domain.HappinessDoorRecord, error) {
	action := result.Actions[0]
	id, _ := strconv.Atoi(action.Identifier)

	var hdr domain.HappinessDoorRecord
	var err error

	switch action.Action {
	case domain.ActionVoteHappy:
		hdr, err = h.repo.IncHappy(id)
	case domain.ActionVoteNeutral:
		hdr, err = h.repo.IncNeutral(id)
	case domain.ActionVoteSad:
		hdr, err = h.repo.IncSad(id)
	}
	return hdr, err
}

func (h *Handlers) Vote(_ http.ResponseWriter, r *http.Request) {
	logRequest(r)
	defer func() { _ = r.Body.Close() }()

	err := r.ParseForm()
	if err != nil {
		log.Println("WARN: Failed to parse form")
		return
	}

	log.Println("Form", r.Form)
	payload, _ := url.QueryUnescape(r.Form.Get("payload"))
	log.Println("Parsed payload", payload)
	var result interactiveResponse

	err = json.Unmarshal([]byte(payload), &result)
	if err != nil {
		log.Println("WARN: Failed to parse response from payload parameter", err)
		return
	}

	hdr, err := h.incrementVoting(result)
	if err != nil {
		log.Println("WARN: Failed to increment voting", err)
	}

	responseUrl := result.ResponseUrl
	log.Println("Got response URL", responseUrl)

	resp := domain.CreateSlackMessage(hdr)
	jsonBytes := toJson(resp)
	_, err = http.Post(responseUrl, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("WARN: Failed to send http request to response URL")
	}

}
