package dtos

import (
	"time"
)

type Week struct {
	Name string `json:"name"`
	Start time.Time `json:"start,omitempty"`
	End time.Time `json:"end,omitempty"`
	Feelings []Feeling `json:"feelings,omitempty"`
}
