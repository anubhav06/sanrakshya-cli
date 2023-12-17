// DO NOT EDIT: generated by syft/internal/packagemetadata/generate/main.go

package packagemetadata

import "github.com/anubhav06/sanrakshya-cli/syft/pkg"

// AllTypes returns a list of all pkg metadata types that syft supports (that are represented in the pkg.Package.Metadata field).
func AllTypes() []any {
	return []any{
		pkg.AlpmDBEntry{},
		pkg.ApkDBEntry{},
		pkg.BinarySignature{},
		pkg.CocoaPodfileLockEntry{},
		pkg.ConanLockEntry{},
		pkg.ConanfileEntry{},
		pkg.ConaninfoEntry{},
		pkg.DartPubspecLockEntry{},
		pkg.DotnetDepsEntry{},
		pkg.DotnetPortableExecutableEntry{},
		pkg.DpkgDBEntry{},
		pkg.ElixirMixLockEntry{},
		pkg.ErlangRebarLockEntry{},
		pkg.GolangBinaryBuildinfoEntry{},
		pkg.GolangModuleEntry{},
		pkg.HackageStackYamlEntry{},
		pkg.HackageStackYamlLockEntry{},
		pkg.JavaArchive{},
		pkg.LinuxKernel{},
		pkg.LinuxKernelModule{},
		pkg.MicrosoftKbPatch{},
		pkg.NixStoreEntry{},
		pkg.NpmPackage{},
		pkg.NpmPackageLockEntry{},
		pkg.PhpComposerInstalledEntry{},
		pkg.PhpComposerLockEntry{},
		pkg.PortageEntry{},
		pkg.PythonPackage{},
		pkg.PythonPipfileLockEntry{},
		pkg.PythonRequirementsEntry{},
		pkg.RDescription{},
		pkg.RpmArchive{},
		pkg.RpmDBEntry{},
		pkg.RubyGemspec{},
		pkg.RustBinaryAuditEntry{},
		pkg.RustCargoLockEntry{},
		pkg.SwiftPackageManagerResolvedEntry{},
	}
}
