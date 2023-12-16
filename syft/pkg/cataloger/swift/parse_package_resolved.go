package swift

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/anchore/syft/internal/log"
	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/pkg/cataloger/generic"
)

var _ generic.Parser = parsePackageResolved

// swift package manager has two versions (1 and 2) of the resolved files, the types below describes the serialization strategies for each version
// with its suffix indicating which version its specific to.

type packageResolvedV1 struct {
	PackageObject packageObjectV1 `json:"object"`
	Version       int             `json:"version"`
}

type packageObjectV1 struct {
	Pins []packagePinsV1
}

type packagePinsV1 struct {
	Name          string       `json:"package"`
	RepositoryURL string       `json:"repositoryURL"`
	State         packageState `json:"state"`
}

type packageResolvedV2 struct {
	Pins []packagePinsV2
}

type packagePinsV2 struct {
	Identity string       `json:"identity"`
	Kind     string       `json:"kind"`
	Location string       `json:"location"`
	State    packageState `json:"state"`
}

type packagePin struct {
	Identity string
	Location string
	Revision string
	Version  string
}

type packageState struct {
	Revision string `json:"revision"`
	Version  string `json:"version"`
}

// parsePackageResolved is a parser for the contents of a Package.resolved file, which is generated by Xcode after it's resolved Swift Package Manger packages.
func parsePackageResolved(_ file.Resolver, _ *generic.Environment, reader file.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	dec := json.NewDecoder(reader)
	var packageResolvedData map[string]interface{}
	for {
		if err := dec.Decode(&packageResolvedData); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, nil, fmt.Errorf("failed to parse Package.resolved file: %w", err)
		}
	}

	if packageResolvedData["version"] == nil {
		log.Trace("no version found in Package.resolved file, skipping")
		return nil, nil, nil
	}

	version, ok := packageResolvedData["version"].(float64)
	if !ok {
		return nil, nil, fmt.Errorf("failed to parse Package.resolved file: version is not a number")
	}

	var pins, err = pinsForVersion(packageResolvedData, version)
	if err != nil {
		return nil, nil, err
	}

	var pkgs []pkg.Package
	for _, pkgPin := range pins {
		pkgs = append(
			pkgs,
			newSwiftPackageManagerPackage(
				pkgPin.Identity,
				pkgPin.Version,
				pkgPin.Location,
				pkgPin.Revision,
				reader.Location.WithAnnotation(pkg.EvidenceAnnotationKey, pkg.PrimaryEvidenceAnnotation),
			),
		)
	}
	return pkgs, nil, nil
}

func pinsForVersion(data map[string]interface{}, version float64) ([]packagePin, error) {
	var genericPins []packagePin
	switch version {
	case 1:
		t := packageResolvedV1{}
		jsonString, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		parseErr := json.Unmarshal(jsonString, &t)
		if parseErr != nil {
			return nil, parseErr
		}
		for _, pin := range t.PackageObject.Pins {
			genericPins = append(genericPins, packagePin{
				pin.Name,
				pin.RepositoryURL,
				pin.State.Revision,
				pin.State.Version,
			})
		}
	case 2:
		t := packageResolvedV2{}
		jsonString, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		parseErr := json.Unmarshal(jsonString, &t)
		if parseErr != nil {
			return nil, parseErr
		}
		for _, pin := range t.Pins {
			genericPins = append(genericPins, packagePin{
				pin.Identity,
				pin.Location,
				pin.State.Revision,
				pin.State.Version,
			})
		}
	default:
		return nil, fmt.Errorf("unknown swift package manager version, %f", version)
	}
	return genericPins, nil
}
