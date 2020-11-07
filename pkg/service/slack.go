package service

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/client"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/db"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
)

type SlackService struct {
	repo        *db.HappinessDoor
	slackClient *client.SlackClient
}

func NewSlackService(repo *db.HappinessDoor, slackClient *client.SlackClient) *SlackService {
	return &SlackService{
		repo:        repo,
		slackClient: slackClient,
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

func (s *SlackService) InitiateHappinessDoor(meetingName, cID, uID string) error {
	member, err := s.isBotMember(cID)
	if err != nil {
		return err
	}
	if !member {
		return s.sendNotAMember(cID, uID)
	}
	return s.sendHappinessDoor(meetingName, cID)
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

	return s.publishVotingToSlack(hdID)
}

func (s *SlackService) InsertUserAction(hdID int, userId string, userName string, action string) error {
	return s.repo.InsertUserAction(hdID, userId, userName, action)
}

func (s *SlackService) sendNotAMember(cID, uID string) error {
	msg := domain.CreateNotAMemberMessage(viper.GetString(conf.BotName))
	return s.slackClient.PostEphemeralMessage(cID, uID, msg)
}

func (s *SlackService) sendHappinessDoor(meetingName string, cID string) error {
	token := generateToken()
	hdID, err := s.repo.CreateHappinessDoor(meetingName, token, cID)
	if err != nil {
		return err
	}

	msg := domain.CreateHappinessDoorMessage(domain.StubRecord(hdID, meetingName))
	msgTS, err := s.slackClient.PostMessage(cID, msg)
	if err != nil {
		return err
	}

	err = s.repo.SetMessageTS(hdID, msgTS)
	return err
}

func (s *SlackService) isBotMember(cID string) (bool, error) {
	return s.slackClient.IsBotMember(cID)
}

func (s *SlackService) publishVotingToSlack(hdID int) error {
	hdr, err := s.computeVoting(hdID)
	if err != nil {
		logrus.WithError(err).WithField("HappinessDoorId", hdID).Warn("Failed to get voting stats")
	}

	msg := domain.CreateHappinessDoorMessage(*hdr)
	return s.slackClient.ReplaceMessage(hdr.ChannelID, hdr.MessageTS, msg)
}

func (s *SlackService) computeVoting(hdID int) (*domain.HappinessDoorDto, error) {
	hdr, err := s.repo.GetHappinessDoorRecord(hdID)
	if err != nil {
		return nil, err
	}

	actions, err := s.repo.GetUserActions(hdID)
	if err != nil {
		return nil, err
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

	return &domain.HappinessDoorDto{
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
