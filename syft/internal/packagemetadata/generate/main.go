package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dave/jennifer/jen"

	"github.com/anubhav06/sanrakshya-cli/syft/internal/packagemetadata"
)

// This program is invoked from syft/internal and generates packagemetadata/generated.go

const (
	pkgImport = "github.com/anubhav06/sanrakshya-cli/syft/pkg"
	path      = "packagemetadata/generated.go"
)

func main() {
	typeNames, err := packagemetadata.DiscoverTypeNames()
	if err != nil {
		panic(fmt.Errorf("unable to get all metadata type names: %w", err))
	}

	// for _, typeName := range typeNames {
	//	fmt.Printf(" - %s\n", typeName)
	//}

	fmt.Printf("updating package metadata type list with %+v types\n", len(typeNames))

	f := jen.NewFile("packagemetadata")
	f.HeaderComment("DO NOT EDIT: generated by syft/internal/packagemetadata/generate/main.go")
	f.ImportName(pkgImport, "pkg")
	f.Comment("AllTypes returns a list of all pkg metadata types that syft supports (that are represented in the pkg.Package.Metadata field).")

	f.Func().Id("AllTypes").Params().Index().Any().BlockFunc(func(g *jen.Group) {
		g.ReturnFunc(func(g *jen.Group) {
			g.Index().Any().ValuesFunc(func(g *jen.Group) {
				for _, typeName := range typeNames {
					g.Qual(pkgImport, typeName).Values()
				}
			})
		})
	})

	rendered := fmt.Sprintf("%#v", f)

	fh, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Errorf("unable to open file: %w", err))
	}

	// fix a little whitespacing
	rendered = strings.ReplaceAll(rendered, ",", ",\n")
	rendered = strings.ReplaceAll(rendered, "[]any{", "[]any{\n")
	rendered = strings.ReplaceAll(rendered, "}}\n}", "},\n}\n}")

	_, err = fh.WriteString(rendered)
	if err != nil {
		panic(fmt.Errorf("unable to write file: %w", err))
	}
	if err := fh.Close(); err != nil {
		panic(fmt.Errorf("unable to close file: %w", err))
	}
}
