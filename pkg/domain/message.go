package domain

import (
	"fmt"
	"github.com/slack-go/slack"
	"strconv"
	"strings"
)

const ActionVoteHappy = "VOTE_HAPPY"
const ActionVoteNeutral = "VOTE_NEUTRAL"
const ActionVoteSad = "VOTE_SAD"

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

func messageWithVotes(message string, votes int) *slack.TextBlockObject {
	m := fmt.Sprintf("%s - *%d votes*", message, votes)
	return markdownText(m)
}

func button(id int, action string, text *slack.TextBlockObject) *slack.ButtonBlockElement {
	return &slack.ButtonBlockElement{
		Type:     "button",
		Text:     text,
		Value:    strconv.Itoa(id),
		ActionID: action,
	}
}

//func greenButton(id int, action string, text *slack.TextBlockObject) *slack.ButtonBlockElement {
//	btn := button(id, action, text)
//	btn.Style = "primary"
//	return btn
//}

func createBlocks(hde HappinessDoorRecord) slack.Blocks {
	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			slack.SectionBlock{
				Type: "section",
				Text: markdownText("How did you find the *" + hde.Name + "* meeting?"),
			},
			slack.NewContextBlock("", plainText("Feel free to leave an additional comment in a thread under this post")),
			slack.DividerBlock{
				Type: "divider",
			},
			slack.SectionBlock{
				Type: "section",
				Text: messageWithVotes(":slightly_smiling_face: I'm happy", hde.Happy),
				Accessory: &slack.Accessory{
					ButtonElement: button(hde.Id, ActionVoteHappy, plainText("Select")),
				},
			},
			slack.SectionBlock{
				Type: "section",
				Text: messageWithVotes(":neutral_face: Neither good nor bad", hde.Neutral),
				Accessory: &slack.Accessory{
					ButtonElement: button(hde.Id, ActionVoteNeutral, plainText("Select")),
				},
			},
			slack.SectionBlock{
				Type: "section",
				Text: messageWithVotes(":disappointed: I did not like it", hde.Sad),
				Accessory: &slack.Accessory{
					ButtonElement: button(hde.Id, ActionVoteSad, plainText("Select")),
				},
			},
			//slack.NewActionBlock("", greenButton(hde.Id, "FEEDBACK", plainText("I want to provide feedback"))), //TODO add feedback
		},
	}
	if len(hde.Voters) > 0 {
		votersString := "*Voted*: " + strings.Join(hde.Voters[:], ", ")

		blocks.BlockSet = append(blocks.BlockSet,
			slack.DividerBlock{
				Type: "divider",
			},
			slack.SectionBlock{
				Type: "section",
				Text: markdownText(votersString),
			})
	}

	return blocks
}

func CreateSlackMessage(hde HappinessDoorRecord) slack.Msg {
	blocks := createBlocks(hde)
	return slack.Msg{Blocks: blocks, ResponseType: "in_channel", ReplaceOriginal: true}
}
