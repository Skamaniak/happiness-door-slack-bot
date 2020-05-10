package domain

import (
	"github.com/slack-go/slack"
	"strconv"
)

func markdownText(text string) *slack.TextBlockObject {
	return &slack.TextBlockObject{
		Type: "mrkdwn",
		Text: text,
	}
}

func plainText(text string) *slack.TextBlockObject {
	return &slack.TextBlockObject{
		Type: "plain_text",
		Text: text,
	}
}

func button(id int, action string, text *slack.TextBlockObject) *slack.ButtonBlockElement {
	return &slack.ButtonBlockElement{
		Type:     "button",
		Text:     text,
		Value:    strconv.Itoa(id),
		ActionID: action,
	}
}

func greenButton(id int, action string, text *slack.TextBlockObject) *slack.ButtonBlockElement {
	btn := button(id, action, text)
	btn.Style = "primary"
	return btn
}

func CreateInitMessage(id int, meetingName string) slack.Msg {
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.SectionBlock{
				Type: "section",
				Text: markdownText("How did you find the *" + meetingName + "* meeting?"),
			},
			slack.DividerBlock{
				Type: "divider",
			},
			slack.SectionBlock{
				Type: "section",
				Text: plainText(":slightly_smiling_face: I'm happy"),
				Accessory: &slack.Accessory{
					ButtonElement: button(id, "VOTE_HAPPY", plainText("Select")),
				},
			},
			slack.SectionBlock{
				Type: "section",
				Text: plainText(":neutral_face: Neither good nor bad"),
				Accessory: &slack.Accessory{
					ButtonElement: button(id, "VOTE_NEUTRAL", plainText("Select")),
				},
			},
			slack.SectionBlock{
				Type: "section",
				Text: plainText(":disappointed: I did not like it"),
				Accessory: &slack.Accessory{
					ButtonElement: button(id, "VOTE_SAD", plainText("Select")),
				},
			},
			slack.DividerBlock{
				Type: "divider",
			},
			slack.NewActionBlock("", greenButton(id, "FEEDBACK", plainText("I want to provide feedback"))),
		},
	}

	return slack.Msg{Blocks: blocks, ResponseType: "in_channel"}
}

func CreateResultMessage() slack.Msg {
	return slack.Msg{Text: "Yay!", ReplaceOriginal: true}
}
