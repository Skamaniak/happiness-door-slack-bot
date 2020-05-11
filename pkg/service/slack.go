package service

import (
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/client"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/db"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	"log"
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

func (s *SlackService) CreateHappinessDoor(meetingName string) (int, error) {
	return s.repo.CreateHappinessDoor(meetingName)
}

func (s *SlackService) GetVoting(hdId int) (*domain.HappinessDoorRecord, error) {
	meetingName, err := s.repo.GetMeetingName(hdId)
	if err != nil {
		return nil, err
	}

	actions, err := s.repo.GetUserActions(hdId)
	if err != nil {
		return nil, err
	}

	var happyVoters, neutralVoters, sadVoters []domain.UserInfo
	for userInfo, action := range actions {
		userIcon, err := s.slackClient.GetUserIconUrl(userInfo.Id)
		if err != nil {
			log.Println("WARN: Failed to fetch icon for user", userInfo.Id, err)
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

	return &domain.HappinessDoorRecord{
		Id:            hdId,
		Name:          meetingName,
		Happy:         len(happyVoters),
		HappyVoters:   happyVoters,
		Neutral:       len(neutralVoters),
		NeutralVoters: neutralVoters,
		Sad:           len(sadVoters),
		SadVoters:     sadVoters,
	}, nil
}

func (s *SlackService) IncrementVoting(result domain.InteractiveResponse) (int, error) {
	action := extractAction(result)
	id := extractHappinessDoorId(result)
	user := result.User

	var err error
	switch action.Action {
	case domain.ActionVoteHappy:
		err = s.repo.InsertUserAction(id, user.Id, user.Name, action.Action)
	case domain.ActionVoteNeutral:
		err = s.repo.InsertUserAction(id, user.Id, user.Name, action.Action)
	case domain.ActionVoteSad:
		err = s.repo.InsertUserAction(id, user.Id, user.Name, action.Action)
	}
	return id, err
}

func (s *SlackService) InsertUserAction(hdId int, userId string, userName string, action string) error {
	return s.repo.InsertUserAction(hdId, userId, userName, action)
}
