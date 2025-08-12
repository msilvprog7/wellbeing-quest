package models

import "time"

type Week struct {
	Id int
	Name string
	Start time.Time
	End time.Time
}
