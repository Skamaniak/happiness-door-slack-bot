package domain

import (
	"fmt"
	"github.com/slack-go/slack"
	"strconv"
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

func image(imageUrl string, altText string) *slack.ImageBlockElement {
	return &slack.ImageBlockElement{
		Type:     "image",
		ImageURL: imageUrl,
		AltText:  altText,
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

func context(element ...slack.MixedElement) *slack.ContextBlock {
	return slack.NewContextBlock("", element...)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func splitIntoChunks(elems []slack.MixedElement, chunkSize int) [][]slack.MixedElement {
	var chunks [][]slack.MixedElement
	for i := 0; i < len(elems); i += chunkSize {
		chunks = append(chunks, elems[i:min(i+chunkSize, len(elems))])
	}
	return chunks
}

func createVoterProfiles(voters []UserVotingAction) []slack.MixedElement {
	var userElems []slack.MixedElement
	for _, userInfo := range voters {
		var userElem slack.MixedElement
		userName := userInfo.Name
		if len(userName) == 0 {
			userName = "-"
		}

		if userInfo.ProfilePicture != "" {
			userElem = image(userInfo.ProfilePicture, userName)
		} else {
			userElem = plainText(userName)
		}
		userElems = append(userElems, userElem)
	}
	return userElems
}

func appendVoters(blockSet []slack.Block, voters []UserVotingAction) []slack.Block {
	userElems := createVoterProfiles(voters)

	if len(userElems) > 0 {
		chunks := splitIntoChunks(userElems, 9)
		for i, chunk := range chunks {
			if i == len(chunks)-1 {
				votes := fmt.Sprintf(" - *%d votes*", len(voters))
				chunk = append(chunk, markdownText(votes))
			}
			blockSet = append(blockSet, context(chunk...))
		}
	}

	return blockSet
}

func CreateHappinessDoorContent(hde HappinessDoorDto) slack.Blocks {
	var blockSet []slack.Block

	blockSet = append(blockSet,
		slack.SectionBlock{
			Type: "section",
			Text: markdownText("How did you find the *" + hde.Name + "* meeting?"),
		},
		slack.DividerBlock{
			Type: "divider",
		},
		slack.SectionBlock{
			Type: "section",
			Text: markdownText(":slightly_smiling_face: I'm happy"),
			Accessory: &slack.Accessory{
				ButtonElement: button(hde.Id, ActionVoteHappy, plainText("Select")),
			},
		},
	)
	blockSet = appendVoters(blockSet, hde.HappyVoters)

	blockSet = append(blockSet, slack.SectionBlock{
		Type: "section",
		Text: markdownText(":neutral_face: Neither good nor bad"),
		Accessory: &slack.Accessory{
			ButtonElement: button(hde.Id, ActionVoteNeutral, plainText("Select")),
		},
	})
	blockSet = appendVoters(blockSet, hde.NeutralVoters)

	blockSet = append(blockSet, slack.SectionBlock{
		Type: "section",
		Text: markdownText(":disappointed: I did not like it"),
		Accessory: &slack.Accessory{
			ButtonElement: button(hde.Id, ActionVoteSad, plainText("Select")),
		},
	})
	blockSet = appendVoters(blockSet, hde.SadVoters)
	blockSet = append(blockSet,
		slack.DividerBlock{
			Type: "divider",
		},
		context(
			markdownText(fmt.Sprintf("Share this happiness door via email using a <%s|web link>.", hde.WebLink)),
			markdownText("Feel free to leave an additional comment in a thread under this post."),
		),
	)
	return slack.Blocks{BlockSet: blockSet}
}

func CreateNotAMemberMessage(botName string) slack.Msg {
	var blockSet []slack.Block

	blockSet = append(blockSet,
		slack.SectionBlock{
			Type: "section",
			Text: markdownText("Bot is not a member of this channel, please execute `/invite @" + botName + "` and try again."),
		},
	)

	return slack.Msg{Blocks: slack.Blocks{BlockSet: blockSet}, ResponseType: slack.ResponseTypeEphemeral, ReplaceOriginal: true}
}
