package models

import "time"

type Week struct {
	Id int
	UserProvider string
	UserId string
	Name string
	Start time.Time
	End time.Time
}
