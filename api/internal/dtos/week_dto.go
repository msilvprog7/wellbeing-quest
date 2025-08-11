package dtos

import (
	"time"
)

type Week struct {
	Name string `json:"name"`
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
	Feelings []Feeling `json:"feelings"`
}
