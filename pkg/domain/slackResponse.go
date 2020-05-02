package domain

import (
	"fmt"
	"io/ioutil"
	"log"
)

func CreateInitMessage(meetingName string) string {
	data, err := ioutil.ReadFile("api/initMessageTemplate.json")
	if err != nil {
		log.Println(err)
	}

	return fmt.Sprintf(string(data), meetingName)
}
