package models

import "time"

type Entry struct {
	Id int
	Activity string
	Feelings []string
	Week string
	Created time.Time
}
