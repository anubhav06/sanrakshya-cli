/*
Package r provides a concrete Cataloger implementation relating to packages within the R language ecosystem.
*/
package r

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

// NewPackageCataloger returns a new R cataloger object based on detection of R package DESCRIPTION files.
func NewPackageCataloger() pkg.Cataloger {
	return generic.NewCataloger("r-package-cataloger").
		WithParserByGlobs(parseDescriptionFile, "**/DESCRIPTION")
}
