package dtos

type Feeling struct {
	Name string `json:"name"`
	Activities []Activity `json:"activities"`
}
