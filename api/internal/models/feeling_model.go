package models

type Feeling struct {
	Id int
	UserProvider string
	UserId string
	Name string
	Activities []string
}
