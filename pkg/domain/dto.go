package domain

import "fmt"

type UserVotingAction struct {
	Id             string
	Name           string
	ProfilePicture string
	Action         string
}

type HappinessDoorDto struct {
	Id            int
	Name          string
	ChannelID     string
	MessageTS     string
	Happy         int
	Neutral       int
	Sad           int
	HappyVoters   []UserVotingAction
	NeutralVoters []UserVotingAction
	SadVoters     []UserVotingAction
	WebLink       string
}

type Action struct {
	Identifier string `json:"value"`
	Action     string `json:"action_id"`
}

type User struct {
	Id string `json:"id"`
}

type InteractiveResponse struct {
	User    User     `json:"user"`
	Actions []Action `json:"actions"`
}

type WsAuth struct {
	UserEmail string
	HdID      int
	Token     string
}

func (a WsAuth) String() string {
	return fmt.Sprintf("userEmail: %s, hdId: %d, token: %s", a.UserEmail, a.HdID, a.Token)
}

func StubRecord(id int, name string) HappinessDoorDto {
	return HappinessDoorDto{Id: id, Name: name}
}
