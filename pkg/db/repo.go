package db

import (
	"database/sql"
	"fmt"
	"github.com/Skamaniak/happiness-door-slack-bot/pkg/conf"
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

func (hd *HappinessDoor) CreateHappinessDoor() (int, error) {
	var id int
	sqlStatement := `INSERT INTO happiness_door DEFAULT VALUES RETURNING id;`
	err := hd.db.QueryRow(sqlStatement).Scan(&id)
	return id, err
}

func (hd *HappinessDoor) increment(id int, attr string) error {
	sqlStatement := "UPDATE happiness_door SET " + attr + " = " + attr + " + 1 WHERE id = $1;"
	_, err := hd.db.Exec(sqlStatement, id)
	return err
}

func (hd *HappinessDoor) IncHappy(id int) error {
	return hd.increment(id, "happy")
}

func (hd *HappinessDoor) IncNeutral(id int) error {
	return hd.increment(id, "neutral")
}

func (hd *HappinessDoor) IncSad(id int) error {
	return hd.increment(id, "sad")
}
