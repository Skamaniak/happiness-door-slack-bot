package domain

type HappinessDoorRecord struct {
	Id      int
	Name    string
	Happy   int
	Neutral int
	Sad     int
}

func StubRecord(id int, name string) HappinessDoorRecord {
	return FullRecord(id, name, 0, 0, 0)
}

func FullRecord(id int, name string, happy int, neutral int, sad int) HappinessDoorRecord {
	return HappinessDoorRecord{
		Id:      id,
		Name:    name,
		Happy:   happy,
		Neutral: neutral,
		Sad:     sad,
	}
}
