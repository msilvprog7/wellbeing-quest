package models

type Activity struct {
	Id int
	UserProvider string
	UserId string
	Name string
	Feelings []string
}
