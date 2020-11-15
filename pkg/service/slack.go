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
	"net/url"
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

func extractHappinessDoorId(action domain.Action) int {
	id := action.Identifier
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

func (s *SlackService) SlackVoting(result domain.InteractiveResponse) error {
	a := extractAction(result)
	u, err := s.slackClient.GetUserById(result.User.Id)

	if err != nil {
		return err
	}

	return s.vote(u, a)
}

func (s *SlackService) WebVoting(ue string, a domain.Action) error {
	u, err := s.slackClient.GetUserById(ue)

	if err != nil {
		return err
	}

	return s.vote(u, a)
}

func (s *SlackService) vote(user client.SlackUser, a domain.Action) error {
	hdID := extractHappinessDoorId(a)

	var err error
	switch a.Action {
	case domain.ActionVoteHappy:
		err = s.repo.InsertUserAction(hdID, user.UID, user.Name, a.Action)
	case domain.ActionVoteNeutral:
		err = s.repo.InsertUserAction(hdID, user.UID, user.Name, a.Action)
	case domain.ActionVoteSad:
		err = s.repo.InsertUserAction(hdID, user.UID, user.Name, a.Action)
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

	var happyVoters, neutralVoters, sadVoters []domain.UserVotingAction
	for _, action := range actions {
		user, err := s.slackClient.GetUserById(action.Id)
		if err != nil {
			logrus.WithError(err).WithField("UserId", action.Id).Warn("Failed to fetch icon for user")
		}
		action.ProfilePicture = user.IconUrl
		switch action.Action {
		case domain.ActionVoteHappy:
			happyVoters = append(happyVoters, action)
		case domain.ActionVoteNeutral:
			neutralVoters = append(neutralVoters, action)
		case domain.ActionVoteSad:
			sadVoters = append(sadVoters, action)
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
		WebLink:       createWebUrl(hdr),
	}, nil
}

func (s *SlackService) VerifyToken(hdId int, token string) (bool, error) {
	return s.repo.HappinessDoorExists(hdId, token)
}

func (s *SlackService) VerifySlackUser(email string) bool {
	_, err := s.slackClient.GetUserByEmail(email) //TODO handle error better (Slack call failure is taken as unverified user)
	return err == nil
}

func (s *SlackService) SubscribeHappinessDoorFeed(hdID int) <-chan domain.HappinessDoorDto {
	return s.pubSub.subscribe(hdID)
}

func (s *SlackService) UnsubscribeHappinessDoorFeed(hdID int, ch <-chan domain.HappinessDoorDto) {
	s.pubSub.unsubscribe(hdID, ch)
}

func createWebUrl(hdr db.HappinessDoorRecord) string {
	webUrl := url.URL{
		Scheme: viper.GetString(conf.WebScheme),
		Host:   viper.GetString(conf.WebHost),
	}
	q := webUrl.Query()
	q.Add("i", strconv.Itoa(hdr.Id))
	q.Add("t", hdr.Token)
	webUrl.RawQuery = q.Encode()
	return webUrl.String()
}
