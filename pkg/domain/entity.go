package domain

type HappinessDoorRecord struct {
	Id      int
	Name    string
	Happy   int
	Neutral int
	Sad     int
	Voters  []string
}

func StubRecord(id int, name string) HappinessDoorRecord {
	return HappinessDoorRecord{Id: id, Name: name}
}
