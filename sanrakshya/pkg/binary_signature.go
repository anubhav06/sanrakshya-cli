package pkg

import "github.com/anubhav06/sanrakshya-cli/sanrakshya/file"

// BinarySignature represents a set of matched values within a binary file.
type BinarySignature struct {
	Matches []ClassifierMatch `mapstructure:"Matches" json:"matches"`
}

// ClassifierMatch represents a single matched value within a binary file and the "class" name the search pattern represents.
type ClassifierMatch struct {
	Classifier string        `mapstructure:"Classifier" json:"classifier"`
	Location   file.Location `mapstructure:"Location" json:"location"`
}
