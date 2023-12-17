/*
Package swift provides a concrete Cataloger implementation relating to packages within the swift language ecosystem.
*/
package swift

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

func NewSwiftPackageManagerCataloger() pkg.Cataloger {
	return generic.NewCataloger("swift-package-manager-cataloger").
		WithParserByGlobs(parsePackageResolved, "**/Package.resolved", "**/.package.resolved")
}

// NewCocoapodsCataloger returns a new Swift Cocoapods lock file cataloger object.
func NewCocoapodsCataloger() pkg.Cataloger {
	return generic.NewCataloger("cocoapods-cataloger").
		WithParserByGlobs(parsePodfileLock, "**/Podfile.lock")
}
