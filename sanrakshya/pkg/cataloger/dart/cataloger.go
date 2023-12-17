/*
Package dart provides a concrete Cataloger implementations for the Dart language ecosystem.
*/
package dart

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

// NewPubspecLockCataloger returns a new Dartlang cataloger object base on pubspec lock files.
func NewPubspecLockCataloger() pkg.Cataloger {
	return generic.NewCataloger("dart-pubspec-lock-cataloger").
		WithParserByGlobs(parsePubspecLock, "**/pubspec.lock")
}
