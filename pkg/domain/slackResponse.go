package domain

import (
	"github.com/slack-go/slack"
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

func button(action string, text *slack.TextBlockObject) *slack.ButtonBlockElement {
	return &slack.ButtonBlockElement{
		Type:  "button",
		Text:  text,
		Value: action,
	}
}

func greenButton(action string, text *slack.TextBlockObject) *slack.ButtonBlockElement {
	btn := button(action, text)
	btn.Style = "primary"
	return btn
}

func CreateInitMessage(meetingName string) slack.Msg {
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
					ButtonElement: button("VOTE_HAPPY", plainText("Select")),
				},
			},
			slack.SectionBlock{
				Type: "section",
				Text: plainText(":neutral_face: Neither good nor bad"),
				Accessory: &slack.Accessory{
					ButtonElement: button("VOTE_NEUTRAL", plainText("Select")),
				},
			},
			slack.SectionBlock{
				Type: "section",
				Text: plainText(":disappointed: I did not like it"),
				Accessory: &slack.Accessory{
					ButtonElement: button("VOTE_SAD", plainText("Select")),
				},
			},
			slack.DividerBlock{
				Type: "divider",
			},
			slack.NewActionBlock("", greenButton("FEEDBACK", plainText("I want to provide feedback"))),
		},
	}

	return slack.Msg{Blocks: blocks}
}

func CreateResultMessage() slack.Msg {
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.SectionBlock{
				Type: "section",
				Text: markdownText("How did you find the meeting?"),
			},
			slack.DividerBlock{
				Type: "divider",
			},
			slack.SectionBlock{
				Type: "section",
				Text: plainText(":slightly_smiling_face: I'm happy"),
			},
			slack.SectionBlock{
				Type: "section",
				Text: plainText(":neutral_face: Neither good nor bad"),
			},
			slack.SectionBlock{
				Type: "section",
				Text: plainText(":disappointed: I did not like it"),
			},
			slack.DividerBlock{
				Type: "divider",
			},
			slack.NewActionBlock("", greenButton("FEEDBACK", plainText("I want to provide feedback"))),
		},
	}

	return slack.Msg{Blocks: blocks, ReplaceOriginal: true}
}
