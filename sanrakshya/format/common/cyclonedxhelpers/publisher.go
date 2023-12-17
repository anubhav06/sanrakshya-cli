package cyclonedxhelpers

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
)

func encodePublisher(p pkg.Package) string {
	if hasMetadata(p) {
		switch metadata := p.Metadata.(type) {
		case pkg.ApkDBEntry:
			return metadata.Maintainer
		case pkg.RpmDBEntry:
			return metadata.Vendor
		case pkg.DpkgDBEntry:
			return metadata.Maintainer
		}
	}
	return ""
}

func decodePublisher(publisher string, metadata interface{}) {
	switch meta := metadata.(type) {
	case *pkg.ApkDBEntry:
		meta.Maintainer = publisher
	case *pkg.RpmDBEntry:
		meta.Vendor = publisher
	case *pkg.DpkgDBEntry:
		meta.Maintainer = publisher
	}
}
