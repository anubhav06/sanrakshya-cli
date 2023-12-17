package elixir

import (
	"github.com/anchore/packageurl-go"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
)

func newPackage(d pkg.ElixirMixLockEntry, locations ...file.Location) pkg.Package {
	p := pkg.Package{
		Name:      d.Name,
		Version:   d.Version,
		Language:  pkg.Elixir,
		Locations: file.NewLocationSet(locations...),
		PURL:      packageURL(d),
		Type:      pkg.HexPkg,
		Metadata:  d,
	}

	p.SetID()

	return p
}

func packageURL(m pkg.ElixirMixLockEntry) string {
	var qualifiers packageurl.Qualifiers

	return packageurl.NewPackageURL(
		packageurl.TypeHex,
		"",
		m.Name,
		m.Version,
		qualifiers,
		"",
	).ToString()
}
