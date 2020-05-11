package domain

type UserInfo struct {
	Id             string
	Name           string
	ProfilePicture string
}

type HappinessDoorRecord struct {
	Id            int
	Name          string
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
	ResponseUrl string   `json:"response_url"`
	User        User     `json:"user"`
	Actions     []Action `json:"actions"`
}

func StubRecord(id int, name string) HappinessDoorRecord {
	return HappinessDoorRecord{Id: id, Name: name}
}
