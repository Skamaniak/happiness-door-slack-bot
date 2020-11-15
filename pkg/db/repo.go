package db

import (
	"database/sql"
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/url"
	"strings"
)

type HappinessDoor struct {
	db *sql.DB
}

type HappinessDoorRecord struct {
	Id          int
	MeetingName string
	ChannelID   string
	MessageTS   string
	Token       string
}

func NewHappinessDoor() (*HappinessDoor, error) {
	db, err := openDb()
	if err != nil {
		return nil, err
	}
	return &HappinessDoor{db: db}, nil
}

func openDb() (*sql.DB, error) {
	dbUrl, err := url.Parse(viper.GetString(conf.DbUrl))
	if err != nil {
		return nil, err
	}

	host := dbUrl.Hostname()
	port := dbUrl.Port()
	user := dbUrl.User.Username()
	password, _ := dbUrl.User.Password()
	dbname := strings.Trim(dbUrl.Path, "/")

	logrus.WithFields(logrus.Fields{
		"DbName": dbname,
		"User":   user,
		"Host":   host + ":" + port,
	}).Info("Connecting to db")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	logrus.Info("Successfully connected to DB")
	return db, nil
}

func (hd *HappinessDoor) CreateHappinessDoor(name string, token string, cID string) (int, error) {
	sqlStatement := `INSERT INTO happiness_door(name, token, channel_id) VALUES ($1, $2, $3) RETURNING id;`
	var id int
	err := hd.db.QueryRow(sqlStatement, name, token, cID).Scan(&id)
	return id, err
}

func (hd *HappinessDoor) SetMessageTS(hdID int, msgTS string) error {
	sqlStatement := `UPDATE happiness_door SET message_ts = $1 WHERE id = $2;`
	_, err := hd.db.Exec(sqlStatement, msgTS, hdID)
	return err
}

func (hd *HappinessDoor) InsertUserAction(hdID int, userId string, userName string, action string) error {
	sqlStatement :=
		`INSERT INTO happiness_door_user_action(happiness_door_id, user_id, user_name, action_id) VALUES ($1, $2, $3, $4)
			ON CONFLICT ON CONSTRAINT unique_user_vote DO UPDATE SET action_id = $4;`

	_, err := hd.db.Exec(sqlStatement, hdID, userId, userName, action)
	return err
}

func (hd *HappinessDoor) GetHappinessDoorRecord(hdID int) (HappinessDoorRecord, error) {
	var r HappinessDoorRecord
	err := hd.db.QueryRow("SELECT id, name, token, channel_id, message_ts FROM happiness_door WHERE id = $1;", hdID).
		Scan(&r.Id, &r.MeetingName, &r.Token, &r.ChannelID, &r.MessageTS)
	return r, err
}

func (hd *HappinessDoor) HappinessDoorExists(hdId int, token string) (bool, error) {
	var count int
	err := hd.db.QueryRow("SELECT count(1) FROM happiness_door WHERE id = $1 AND token = $2", hdId, token).Scan(&count)
	return count != 0, err
}

func (hd *HappinessDoor) GetUserActions(hdID int) ([]domain.UserVotingAction, error) {
	var actions []domain.UserVotingAction

	rows, err := hd.db.Query(
		"SELECT action_id, user_id, user_name from happiness_door_user_action WHERE happiness_door_id = $1 ORDER BY user_name ASC;", hdID)
	if err != nil {
		return actions, err
	}

	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var userAction domain.UserVotingAction
		err := rows.Scan(&userAction.Action, &userAction.Id, &userAction.Name)
		if err != nil {
			return actions, err
		}
		actions = append(actions, userAction)
	}
	err = rows.Err()
	return actions, err
}
