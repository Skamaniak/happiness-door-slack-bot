package service

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/client"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/db"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
	"strconv"
)

type SlackService struct {
	repo        *db.HappinessDoor
	slackClient *client.SlackClient
	pubSub      *happinessDoorPubSub
}

func NewSlackService(repo *db.HappinessDoor, slackClient *client.SlackClient) *SlackService {
	return &SlackService{
		repo:        repo,
		slackClient: slackClient,
		pubSub:      newPubsub(),
	}
}

func extractAction(res domain.InteractiveResponse) domain.Action {
	return res.Actions[0]
}

func extractHappinessDoorId(res domain.InteractiveResponse) int {
	id := extractAction(res).Identifier
	i, _ := strconv.Atoi(id)
	return i
}

func generateToken() string {
	tl := viper.GetInt(conf.WebTokenLength)
	b := make([]byte, tl/2)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func (s *SlackService) InitiateHappinessDoor(meetingName, cID string) (*slack.Msg, error) {
	canPost, err := s.canPostMessageToChannel(cID)
	if err != nil {
		return nil, err
	}
	if !canPost {
		msg := domain.CreateNotAMemberMessage(viper.GetString(conf.BotName))
		return &msg, nil
	}
	return nil, s.sendHappinessDoor(meetingName, cID)
}

func (s *SlackService) IncrementVoting(result domain.InteractiveResponse) error {
	action := extractAction(result)
	hdID := extractHappinessDoorId(result)
	user := result.User

	var err error
	switch action.Action {
	case domain.ActionVoteHappy:
		err = s.repo.InsertUserAction(hdID, user.Id, user.Name, action.Action)
	case domain.ActionVoteNeutral:
		err = s.repo.InsertUserAction(hdID, user.Id, user.Name, action.Action)
	case domain.ActionVoteSad:
		err = s.repo.InsertUserAction(hdID, user.Id, user.Name, action.Action)
	}
	if err != nil {
		return err
	}

	return s.publishVoting(hdID)
}

func (s *SlackService) sendHappinessDoor(meetingName string, cID string) error {
	token := generateToken()
	hdID, err := s.repo.CreateHappinessDoor(meetingName, token, cID)
	if err != nil {
		return err
	}

	msg := domain.CreateHappinessDoorContent(domain.StubRecord(hdID, meetingName))
	msgTS, err := s.slackClient.PostMessage(cID, msg)
	if err != nil {
		return err
	}

	err = s.repo.SetMessageTS(hdID, msgTS)
	return err
}

func (s *SlackService) canPostMessageToChannel(cID string) (bool, error) {
	return s.slackClient.CanPostMessage(cID)
}

func (s *SlackService) publishVoting(hdID int) error {
	hdr, err := s.ComputeVoting(hdID)
	if err != nil {
		logrus.WithError(err).WithField("HappinessDoorId", hdID).Warn("Failed to get voting stats")
	}

	s.pubSub.publish(hdr)
	msg := domain.CreateHappinessDoorContent(hdr)
	return s.slackClient.ReplaceMessage(hdr.ChannelID, hdr.MessageTS, msg)
}

func (s *SlackService) ComputeVoting(hdID int) (domain.HappinessDoorDto, error) {
	hdr, err := s.repo.GetHappinessDoorRecord(hdID)
	if err != nil {
		return domain.HappinessDoorDto{}, err
	}

	actions, err := s.repo.GetUserActions(hdID)
	if err != nil {
		return domain.HappinessDoorDto{}, err
	}

	var happyVoters, neutralVoters, sadVoters []domain.UserInfo
	for userInfo, action := range actions {
		userIcon, err := s.slackClient.GetUserIconUrl(userInfo.Id)
		if err != nil {
			logrus.WithError(err).WithField("UserId", userInfo.Id).Warn("Failed to fetch icon for user")
		}
		userInfo.ProfilePicture = userIcon
		switch action {
		case domain.ActionVoteHappy:
			happyVoters = append(happyVoters, userInfo)
		case domain.ActionVoteNeutral:
			neutralVoters = append(neutralVoters, userInfo)
		case domain.ActionVoteSad:
			sadVoters = append(sadVoters, userInfo)
		}
	}

	return domain.HappinessDoorDto{
		Id:            hdID,
		Name:          hdr.MeetingName,
		ChannelID:     hdr.ChannelID,
		MessageTS:     hdr.MessageTS,
		Happy:         len(happyVoters),
		HappyVoters:   happyVoters,
		Neutral:       len(neutralVoters),
		NeutralVoters: neutralVoters,
		Sad:           len(sadVoters),
		SadVoters:     sadVoters,
	}, nil
}

func (s *SlackService) VerifyToken(hdId, token string) (bool, error) {
	return s.repo.HappinessDoorExists(hdId, token)
}

func (s *SlackService) SubscribeHappinessDoorFeed(hdID int) <-chan domain.HappinessDoorDto {
	return s.pubSub.subscribe(hdID)
}

func (s *SlackService) UnsubscribeHappinessDoorFeed(hdID int, ch <-chan domain.HappinessDoorDto) {
	s.pubSub.unsubscribe(hdID, ch)
}
