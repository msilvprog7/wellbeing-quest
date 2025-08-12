package dtos

import (
	"time"
)

type Activity struct {
	Name string `json:"name"`
	Feelings []string `json:"feelings,omitempty"`
	Created time.Time `json:"created,omitempty"`
	Week string `json:"week,omitempty"`
}
