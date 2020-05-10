package db

import (
	"database/sql"
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/domain"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"strconv"
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
	host := viper.GetString(conf.DbHost)
	port := viper.GetInt(conf.DbPort)
	user := viper.GetString(conf.DbUser)
	password := viper.GetString(conf.DbPassword)
	dbname := viper.GetString(conf.DbName)

	log.Println("INFO: Connecting to db", dbname, "on host", host+":"+strconv.Itoa(port))

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
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

func (hd *HappinessDoor) InsertUserAction(hdId int, userId string, action string) error {
	sqlStatement := `INSERT INTO happiness_door_user_action(happiness_door_id, user_id, action_id) VALUES ($1, $2, $3)
			ON CONFLICT ON CONSTRAINT unique_user_vote DO UPDATE SET action_id = $3;`

	_, err := hd.db.Exec(sqlStatement, hdId, userId, action)
	return err
}

// TODO create a service layer and move this to it
func (hd *HappinessDoor) GetStats(hdId int) (*domain.HappinessDoorRecord, error) {
	var r domain.HappinessDoorRecord
	err := hd.db.QueryRow("SELECT id, name FROM happiness_door WHERE id = $1;", hdId).Scan(&r.Id, &r.Name)

	if err != nil {
		return nil, err
	}

	rows, err := hd.db.Query("SELECT action_id from happiness_door_user_action WHERE happiness_door_id = $1;", hdId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var action string
		err := rows.Scan(&action)
		if err != nil {
			return nil, err
		}
		switch action {
		case domain.ActionVoteHappy:
			r.Happy++
		case domain.ActionVoteNeutral:
			r.Neutral++
		case domain.ActionVoteSad:
			r.Sad++
		}
	}

	err = rows.Err()
	return &r, err
}
