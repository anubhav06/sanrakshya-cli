package spdxhelpers

import "github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"

func Description(p pkg.Package) string {
	if hasMetadata(p) {
		switch metadata := p.Metadata.(type) {
		case pkg.ApkDBEntry:
			return metadata.Description
		case pkg.NpmPackage:
			return metadata.Description
		}
	}
	return ""
}

func hasMetadata(p pkg.Package) bool {
	return p.Metadata != nil
}
