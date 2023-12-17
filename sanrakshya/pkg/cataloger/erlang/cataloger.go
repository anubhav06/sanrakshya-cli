/*
Package erlang provides a concrete Cataloger implementation relating to packages within the Erlang language ecosystem.
*/
package erlang

import (
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg/cataloger/generic"
)

// NewRebarLockCataloger returns a new cataloger instance for Erlang rebar.lock files.
func NewRebarLockCataloger() pkg.Cataloger {
	return generic.NewCataloger("erlang-rebar-lock-cataloger").
		WithParserByGlobs(parseRebarLock, "**/rebar.lock")
}
