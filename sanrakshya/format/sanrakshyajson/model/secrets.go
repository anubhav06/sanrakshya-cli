package model

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file"
)

type Secrets struct {
	Location file.Coordinates    `json:"location"`
	Secrets  []file.SearchResult `json:"secrets"`
}
