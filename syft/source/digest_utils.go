package source

import (
	"strings"

	"github.com/anubhav06/sanrakshya-cli/syft/artifact"
)

func artifactIDFromDigest(input string) artifact.ID {
	return artifact.ID(strings.TrimPrefix(input, "sha256:"))
}
