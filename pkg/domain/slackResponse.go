package domain

import (
	"github.com/slack-go/slack"
)

func CreateInitMessage(meetingName string) slack.Msg {
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.SectionBlock{
				Type: "section",
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "How did you find the *" + meetingName + "* meeting?",
				},
			},
			slack.DividerBlock{
				Type: "divider",
			},
			slack.SectionBlock{
				Type: "section",
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: ":slightly_smiling_face: I'm happy",
				},
				Accessory: &slack.Accessory{
					ButtonElement: &slack.ButtonBlockElement{
						Type: "button",
						Text: &slack.TextBlockObject{
							Type: "plain_text",
							Text: "Select",
						},
						Value: "VOTE_HAPPY",
					},
				},
			},
			slack.SectionBlock{
				Type: "section",
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: ":neutral_face: Neither good nor bad",
				},
				Accessory: &slack.Accessory{
					ButtonElement: &slack.ButtonBlockElement{
						Type: "button",
						Text: &slack.TextBlockObject{
							Type: "plain_text",
							Text: "Select",
						},
						Value: "VOTE_NEUTRAL",
					},
				},
			},
			slack.SectionBlock{
				Type: "section",
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: ":disappointed: I did not like it",
				},
				Accessory: &slack.Accessory{
					ButtonElement: &slack.ButtonBlockElement{
						Type: "button",
						Text: &slack.TextBlockObject{
							Type: "plain_text",
							Text: "Select",
						},
						Value: "VOTE_SAD",
					},
				},
			},
			slack.DividerBlock{
				Type: "divider",
			},
			slack.ActionBlock{
				Type: "actions",
				Elements: slack.BlockElements{
					ElementSet: []slack.BlockElement{
						slack.ButtonBlockElement{
							Text: &slack.TextBlockObject{
								Type: "plain_text",
								Text: "I want to provide feedback",
							},
							Style: "primary",
							Value: "VOTE_SAD",
						},
					},
				},
			},
		},
	}

	return slack.Msg{Blocks: blocks}
}
