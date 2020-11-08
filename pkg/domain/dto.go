package domain

type UserInfo struct {
	Id             string
	Name           string
	ProfilePicture string
}

type HappinessDoorDto struct {
	Id            int
	Name          string
	ChannelID     string
	MessageTS     string
	Happy         int
	Neutral       int
	Sad           int
	HappyVoters   []UserInfo
	NeutralVoters []UserInfo
	SadVoters     []UserInfo
}

type Action struct {
	Identifier string `json:"value"`
	Action     string `json:"action_id"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type InteractiveResponse struct {
	User    User     `json:"user"`
	Actions []Action `json:"actions"`
}

func StubRecord(id int, name string) HappinessDoorDto {
	return HappinessDoorDto{Id: id, Name: name}
}
