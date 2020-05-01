package domain

type SlackResponse struct {
	Markdown bool   `json:"mrkdwn"`
	Text     string `json:"text"`
}
