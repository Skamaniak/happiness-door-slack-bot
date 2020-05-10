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

func (hd *HappinessDoor) increment(id int, attr string) (domain.HappinessDoorRecord, error) {
	sqlStatement := "UPDATE happiness_door SET " + attr + " = " + attr + " + 1 WHERE id = $1 RETURNING id, name, happy, neutral, sad;"
	var r domain.HappinessDoorRecord

	err := hd.db.QueryRow(sqlStatement, id).Scan(&r.Id, &r.Name, &r.Happy, &r.Neutral, &r.Sad)
	return r, err
}

func (hd *HappinessDoor) IncHappy(id int) (domain.HappinessDoorRecord, error) {
	return hd.increment(id, "happy")
}

func (hd *HappinessDoor) IncNeutral(id int) (domain.HappinessDoorRecord, error) {
	return hd.increment(id, "neutral")
}

func (hd *HappinessDoor) IncSad(id int) (domain.HappinessDoorRecord, error) {
	return hd.increment(id, "sad")
}
