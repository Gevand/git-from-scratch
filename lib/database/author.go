package database

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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
	temp_split := strings.Split(line, "> ")
	if len(temp_split) != 2 {
		return errors.New("author is not in a proper format")
	}
	timestamp, err := strconv.ParseInt(strings.Replace(temp_split[1], " +0000", "", 1), 10, 64)
	if err != nil {
		return err
	}
	a.Time = time.Unix(timestamp, 0)
	temp_split = strings.Split(temp_split[0], " <")
	if len(temp_split) != 2 {
		return errors.New("author is not in a proper format")
	}
	a.Name = temp_split[0]
	a.Email = temp_split[1]
	return err
}
