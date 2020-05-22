package db

import (
	"database/sql"
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"strings"
)

type HappinessDoor struct {
	db *sql.DB
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

	log.Println("INFO: Connecting to db", dbname, "on host", host+":"+port)
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

	log.Println("INFO: Successfully connected to DB")
	return db, nil
}

func (hd *HappinessDoor) CreateHappinessDoor(name string) (int, error) {
	sqlStatement := `INSERT INTO happiness_door(name) VALUES ($1) RETURNING id;`
	var id int
	err := hd.db.QueryRow(sqlStatement, name).Scan(&id)
	return id, err
}

func (hd *HappinessDoor) InsertUserAction(hdId int, userId string, userName string, action string) error {
	sqlStatement :=
		`INSERT INTO happiness_door_user_action(happiness_door_id, user_id, user_name, action_id) VALUES ($1, $2, $3, $4)
			ON CONFLICT ON CONSTRAINT unique_user_vote DO UPDATE SET action_id = $4;`

	_, err := hd.db.Exec(sqlStatement, hdId, userId, userName, action)
	return err
}

func (hd *HappinessDoor) GetMeetingName(hdId int) (string, error) {
	var name string
	err := hd.db.QueryRow("SELECT name FROM happiness_door WHERE id = $1;", hdId).Scan(&name)
	return name, err
}

func (hd *HappinessDoor) GetUserActions(hdId int) (map[domain.UserInfo]string, error) {
	actions := make(map[domain.UserInfo]string)

	rows, err := hd.db.Query(
		"SELECT action_id, user_id, user_name from happiness_door_user_action WHERE happiness_door_id = $1 ORDER BY user_name ASC;", hdId)
	if err != nil {
		return actions, err
	}

	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var userInfo domain.UserInfo
		var action string
		err := rows.Scan(&action, &userInfo.Id, &userInfo.Name)
		if err != nil {
			return actions, err
		}
		actions[userInfo] = action
	}
	err = rows.Err()
	return actions, err
}
