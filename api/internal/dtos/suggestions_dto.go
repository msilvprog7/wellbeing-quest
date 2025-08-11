package dtos

type Suggestions struct {
	Activities []Activity `json:"activities"`
	Feelings []Feeling `json:"feelings"`
}
