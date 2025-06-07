package database

import (
	"fmt"
	"time"
)

type Author struct {
	Name, Email string
	Time        time.Time
}

func NewAuthor(name, email string, time time.Time) *Author {
	return &Author{Name: name, Email: email, Time: time}
}

func (a *Author) ToString() string {
	return_string := fmt.Sprintf("%s <%s> %v +0000", a.Name, a.Email, a.Time.Unix())
	return return_string
}

func (a *Author) Parse(line string) error {
	var timestamp int64
	_, err := fmt.Sscanf(line, "%s <%s> %v +0000", a.Name, a.Email, timestamp)
	a.Time = time.Unix(timestamp, 0)
	return err
}
