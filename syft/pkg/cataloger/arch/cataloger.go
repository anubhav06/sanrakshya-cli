/*
Package arch provides a concrete Cataloger implementations for packages relating to the Arch linux distribution.
*/
package arch

import (
	"github.com/anubhav06/sanrakshya-cli/syft/pkg"
	"github.com/anubhav06/sanrakshya-cli/syft/pkg/cataloger/generic"
)

// NewDBCataloger returns a new cataloger object initialized for arch linux pacman database flat-file stores.
func NewDBCataloger() pkg.Cataloger {
	return generic.NewCataloger("alpm-db-cataloger").
		WithParserByGlobs(parseAlpmDB, pkg.AlpmDBGlob)
}
