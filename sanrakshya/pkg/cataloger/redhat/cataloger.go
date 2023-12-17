/*
Package redhat provides a concrete DBCataloger implementation relating to packages within the RedHat linux distribution.
*/
package redhat

import (
	"database/sql"

	"github.com/anubhav06/sanrakshya-cli/internal/log"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

// NewDBCataloger returns a new RPM DB cataloger object.
func NewDBCataloger() pkg.Cataloger {
	// check if a sqlite driver is available
	if !isSqliteDriverAvailable() {
		log.Warnf("sqlite driver is not available, newer RPM databases might not be cataloged")
	}

	return generic.NewCataloger("rpm-db-cataloger").
		WithParserByGlobs(parseRpmDB, pkg.RpmDBGlob).
		WithParserByGlobs(parseRpmManifest, pkg.RpmManifestGlob)
}

// NewArchiveCataloger returns a new RPM file cataloger object.
func NewArchiveCataloger() pkg.Cataloger {
	return generic.NewCataloger("rpm-archive-cataloger").
		WithParserByGlobs(parseRpmArchive, "**/*.rpm")
}

func isSqliteDriverAvailable() bool {
	_, err := sql.Open("sqlite", ":memory:")
	return err == nil
}
