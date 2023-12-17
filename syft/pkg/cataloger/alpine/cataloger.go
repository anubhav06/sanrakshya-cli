/*
Package alpine provides a concrete Cataloger implementations for packages relating to the Alpine linux distribution.
*/
package alpine

import (
	"github.com/anubhav06/sanrakshya-cli/syft/pkg"
	"github.com/anubhav06/sanrakshya-cli/syft/pkg/cataloger/generic"
)

// NewDBCataloger returns a new cataloger object initialized for Alpine package DB flat-file stores.
func NewDBCataloger() pkg.Cataloger {
	return generic.NewCataloger("apk-db-cataloger").
		WithParserByGlobs(parseApkDB, pkg.ApkDBGlob)
}
