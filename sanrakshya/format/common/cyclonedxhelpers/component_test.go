package cyclonedxhelpers

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/CycloneDX/cyclonedx-go"
	"github.com/stretchr/testify/assert"

	"github.com/anubhav06/sanrakshya-cli/sanrakshya/file"
	"github.com/anubhav06/sanrakshya-cli/sanrakshya/pkg"
)

func Test_encodeComponentProperties(t *testing.T) {
	epoch := 2
	tests := []struct {
		name     string
		input    pkg.Package
		expected []cyclonedx.Property
	}{
		{
			name:     "no metadata",
			input:    pkg.Package{},
			expected: nil,
		},
		{
			name: "from apk",
			input: pkg.Package{
				FoundBy: "cataloger",
				Locations: file.NewLocationSet(
					file.NewLocationFromCoordinates(file.Coordinates{RealPath: "test"}),
				),
				Metadata: pkg.ApkDBEntry{
					Package:       "libc-utils",
					OriginPackage: "libc-dev",
					Maintainer:    "Natanael Copa <ncopa@alpinelinux.org>",
					Version:       "0.7.2-r0",
					Architecture:  "x86_64",
					URL:           "http://alpinelinux.org",
					Description:   "Meta package to pull in correct libc",
					Size:          0,
					InstalledSize: 4096,
					Dependencies:  []string{"musl-utils"},
					Provides:      []string{"so:libc.so.1"},
					Checksum:      "Q1p78yvTLG094tHE1+dToJGbmYzQE=",
					GitCommit:     "97b1c2842faa3bfa30f5811ffbf16d5ff9f1a479",
					Files:         []pkg.ApkFileRecord{},
				},
			},
			expected: []cyclonedx.Property{
				{Name: "sanrakshya:package:foundBy", Value: "cataloger"},
				{Name: "sanrakshya:package:metadataType", Value: "apk-db-entry"},
				{Name: "sanrakshya:location:0:path", Value: "test"},
				{Name: "sanrakshya:metadata:gitCommitOfApkPort", Value: "97b1c2842faa3bfa30f5811ffbf16d5ff9f1a479"},
				{Name: "sanrakshya:metadata:installedSize", Value: "4096"},
				{Name: "sanrakshya:metadata:originPackage", Value: "libc-dev"},
				{Name: "sanrakshya:metadata:provides:0", Value: "so:libc.so.1"},
				{Name: "sanrakshya:metadata:pullChecksum", Value: "Q1p78yvTLG094tHE1+dToJGbmYzQE="},
				{Name: "sanrakshya:metadata:pullDependencies:0", Value: "musl-utils"},
				{Name: "sanrakshya:metadata:size", Value: "0"},
			},
		},
		{
			name: "from dpkg",
			input: pkg.Package{
				Metadata: pkg.DpkgDBEntry{
					Package:       "tzdata",
					Version:       "2020a-0+deb10u1",
					Source:        "tzdata-dev",
					SourceVersion: "1.0",
					Architecture:  "all",
					InstalledSize: 3036,
					Maintainer:    "GNU Libc Maintainers <debian-glibc@lists.debian.org>",
					Files:         []pkg.DpkgFileRecord{},
				},
			},
			expected: []cyclonedx.Property{
				{Name: "sanrakshya:package:metadataType", Value: "dpkg-db-entry"},
				{Name: "sanrakshya:metadata:installedSize", Value: "3036"},
				{Name: "sanrakshya:metadata:source", Value: "tzdata-dev"},
				{Name: "sanrakshya:metadata:sourceVersion", Value: "1.0"},
			},
		},
		{
			name: "from go bin",
			input: pkg.Package{
				Name:     "golang.org/x/net",
				Version:  "v0.0.0-20211006190231-62292e806868",
				Language: pkg.Go,
				Type:     pkg.GoModulePkg,
				Metadata: pkg.GolangBinaryBuildinfoEntry{
					GoCompiledVersion: "1.17",
					Architecture:      "amd64",
					H1Digest:          "h1:KlOXYy8wQWTUJYFgkUI40Lzr06ofg5IRXUK5C7qZt1k=",
				},
			},
			expected: []cyclonedx.Property{
				{Name: "sanrakshya:package:language", Value: pkg.Go.String()},
				{Name: "sanrakshya:package:metadataType", Value: "go-module-buildinfo-entry"},
				{Name: "sanrakshya:package:type", Value: "go-module"},
				{Name: "sanrakshya:metadata:architecture", Value: "amd64"},
				{Name: "sanrakshya:metadata:goCompiledVersion", Value: "1.17"},
				{Name: "sanrakshya:metadata:h1Digest", Value: "h1:KlOXYy8wQWTUJYFgkUI40Lzr06ofg5IRXUK5C7qZt1k="},
			},
		},
		{
			name: "from go mod",
			input: pkg.Package{
				Name:     "golang.org/x/net",
				Version:  "v0.0.0-20211006190231-62292e806868",
				Language: pkg.Go,
				Type:     pkg.GoModulePkg,
				Metadata: pkg.GolangModuleEntry{
					H1Digest: "h1:KlOXYy8wQWTUJYFgkUI40Lzr06ofg5IRXUK5C7qZt1k=",
				},
			},
			expected: []cyclonedx.Property{
				{Name: "sanrakshya:package:language", Value: pkg.Go.String()},
				{Name: "sanrakshya:package:metadataType", Value: "go-module-entry"},
				{Name: "sanrakshya:package:type", Value: "go-module"},
				{Name: "sanrakshya:metadata:h1Digest", Value: "h1:KlOXYy8wQWTUJYFgkUI40Lzr06ofg5IRXUK5C7qZt1k="},
			},
		},
		{
			name: "from rpm",
			input: pkg.Package{
				Name:    "dive",
				Version: "0.9.2-1",
				Type:    pkg.RpmPkg,
				Metadata: pkg.RpmDBEntry{
					Name:      "dive",
					Epoch:     &epoch,
					Arch:      "x86_64",
					Release:   "1",
					Version:   "0.9.2",
					SourceRpm: "dive-0.9.2-1.src.rpm",
					Size:      12406784,
					Vendor:    "",
					Files:     []pkg.RpmFileRecord{},
				},
			},
			expected: []cyclonedx.Property{
				{Name: "sanrakshya:package:metadataType", Value: "rpm-db-entry"},
				{Name: "sanrakshya:package:type", Value: "rpm"},
				{Name: "sanrakshya:metadata:epoch", Value: "2"},
				{Name: "sanrakshya:metadata:release", Value: "1"},
				{Name: "sanrakshya:metadata:size", Value: "12406784"},
				{Name: "sanrakshya:metadata:sourceRpm", Value: "dive-0.9.2-1.src.rpm"},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := encodeComponent(test.input)
			if test.expected == nil {
				if c.Properties != nil {
					t.Fatalf("expected no properties, got: %+v", *c.Properties)
				}
				return
			}
			assert.ElementsMatch(t, test.expected, *c.Properties)
		})
	}
}

func Test_encodeCompomentType(t *testing.T) {
	tests := []struct {
		name string
		pkg  pkg.Package
		want cyclonedx.Component
	}{
		{
			name: "non-binary package",
			pkg: pkg.Package{
				Name:    "pkg1",
				Version: "1.9.2",
				Type:    pkg.GoModulePkg,
			},
			want: cyclonedx.Component{
				Name:    "pkg1",
				Version: "1.9.2",
				Type:    cyclonedx.ComponentTypeLibrary,
				Properties: &[]cyclonedx.Property{
					{
						Name:  "sanrakshya:package:type",
						Value: "go-module",
					},
				},
			},
		},
		{
			name: "non-binary package",
			pkg: pkg.Package{
				Name:    "pkg1",
				Version: "3.1.2",
				Type:    pkg.BinaryPkg,
			},
			want: cyclonedx.Component{
				Name:    "pkg1",
				Version: "3.1.2",
				Type:    cyclonedx.ComponentTypeApplication,
				Properties: &[]cyclonedx.Property{
					{
						Name:  "sanrakshya:package:type",
						Value: "binary",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pkg.ID()
			p := encodeComponent(tt.pkg)
			assert.Equal(t, tt.want, p)
		})
	}
}

func Test_deriveBomRef(t *testing.T) {
	pkgWithPurl := pkg.Package{
		Name:    "django",
		Version: "1.11.1",
		PURL:    "pkg:pypi/django@1.11.1",
	}
	pkgWithPurl.SetID()

	pkgWithOutPurl := pkg.Package{
		Name:    "django",
		Version: "1.11.1",
		PURL:    "",
	}
	pkgWithOutPurl.SetID()

	pkgWithBadPurl := pkg.Package{
		Name:    "django",
		Version: "1.11.1",
		PURL:    "pkg:pyjango@1.11.1",
	}
	pkgWithBadPurl.SetID()

	tests := []struct {
		name string
		pkg  pkg.Package
		want string
	}{
		{
			name: "use pURL-id hybrid",
			pkg:  pkgWithPurl,
			want: fmt.Sprintf("pkg:pypi/django@1.11.1?package-id=%s", pkgWithPurl.ID()),
		},
		{
			name: "fallback to ID when pURL is invalid",
			pkg:  pkgWithBadPurl,
			want: string(pkgWithBadPurl.ID()),
		},
		{
			name: "fallback to ID when pURL is missing",
			pkg:  pkgWithOutPurl,
			want: string(pkgWithOutPurl.ID()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.pkg.ID()
			assert.Equal(t, tt.want, deriveBomRef(tt.pkg))
		})
	}
}

func Test_decodeComponent(t *testing.T) {
	tests := []struct {
		name         string
		component    cyclonedx.Component
		wantLanguage pkg.Language
		wantMetadata any
	}{
		{
			name: "derive language from pURL if missing",
			component: cyclonedx.Component{
				Name:       "ch.qos.logback/logback-classic",
				Version:    "1.2.3",
				PackageURL: "pkg:maven/ch.qos.logback/logback-classic@1.2.3",
				Type:       "library",
				BOMRef:     "pkg:maven/ch.qos.logback/logback-classic@1.2.3",
			},
			wantLanguage: pkg.Java,
		},
		{
			name: "handle RpmdbMetadata type without properties",
			component: cyclonedx.Component{
				Name:       "acl",
				Version:    "2.2.53-1.el8",
				PackageURL: "pkg:rpm/centos/acl@2.2.53-1.el8?arch=x86_64&upstream=acl-2.2.53-1.el8.src.rpm&distro=centos-8",
				Type:       "library",
				BOMRef:     "pkg:rpm/centos/acl@2.2.53-1.el8?arch=x86_64&upstream=acl-2.2.53-1.el8.src.rpm&distro=centos-8",
				Properties: &[]cyclonedx.Property{
					{
						Name:  "sanrakshya:package:metadataType",
						Value: "RpmdbMetadata",
					},
				},
			},
			wantMetadata: pkg.RpmDBEntry{},
		},
		{
			name: "handle RpmdbMetadata type with properties",
			component: cyclonedx.Component{
				Name:       "acl",
				Version:    "2.2.53-1.el8",
				PackageURL: "pkg:rpm/centos/acl@2.2.53-1.el8?arch=x86_64&upstream=acl-2.2.53-1.el8.src.rpm&distro=centos-8",
				Type:       "library",
				BOMRef:     "pkg:rpm/centos/acl@2.2.53-1.el8?arch=x86_64&upstream=acl-2.2.53-1.el8.src.rpm&distro=centos-8",
				Properties: &[]cyclonedx.Property{
					{
						Name:  "sanrakshya:package:metadataType",
						Value: "RpmDBMetadata",
					},
					{
						Name:  "sanrakshya:metadata:release",
						Value: "some-release",
					},
				},
			},
			wantMetadata: pkg.RpmDBEntry{
				Release: "some-release",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := decodeComponent(&tt.component)
			if tt.wantLanguage != "" {
				assert.Equal(t, tt.wantLanguage, p.Language)
			}
			if tt.wantMetadata != nil {
				assert.Truef(t, reflect.DeepEqual(tt.wantMetadata, p.Metadata), "metadata should match: %+v != %+v", tt.wantMetadata, p.Metadata)
			}
			if tt.wantMetadata == nil && tt.wantLanguage == "" {
				t.Fatal("this is a useless test, please remove it")
			}
		})
	}
}
