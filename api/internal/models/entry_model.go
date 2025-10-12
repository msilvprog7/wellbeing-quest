package models

import "time"

type Entry struct {
	Id int
	UserProvider string
	UserId string
	Activity string
	Feelings []string
	Week string
	Created time.Time
}
