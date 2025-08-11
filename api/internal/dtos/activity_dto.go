package dtos

import (
	"time"
)

type Activity struct {
	Name string `json:"name"`
	Feelings []string `json:"feelings"`
	Created time.Time `json:"created"`
	Week string `json:"week"`
}
