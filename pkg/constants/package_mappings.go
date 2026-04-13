package constants

// PackageMapping represents the mapping between macOS and Debian package names
type PackageMapping struct {
	MacOS  string // Homebrew package name
	Debian string // Debian/Ubuntu package name
}

// PackageMappings maps package constants to their platform-specific names
// This enables seamless translation from macOS Homebrew names to Debian/Ubuntu apt names
var PackageMappings = map[string]PackageMapping{
	Gdbm: {
		MacOS:  "gdbm",
		Debian: "libgdbm-dev",
	},
	Jemalloc: {
		MacOS:  "jemalloc",
		Debian: "libjemalloc2",
	},
	Libffi: {
		MacOS:  "libffi",
		Debian: "libffi-dev",
	},
	Libyaml: {
		MacOS:  "libyaml",
		Debian: "libyaml-dev",
	},
	Ncurses: {
		MacOS:  "ncurses",
		Debian: "libncurses5-dev",
	},
	Readline: {
		MacOS:  "readline",
		Debian: "libreadline-dev",
	},
	Vips: {
		MacOS:  "vips",
		Debian: "libvips",
	},
	Zlib: {
		MacOS:  "zlib",
		Debian: "zlib1g-dev",
	},
}

// GetDebianPackageName returns the Debian package name for a given package constant
// If the package is not found in the mappings, it returns the original name as a fallback
func GetDebianPackageName(packageConstant string) string {
	if mapping, exists := PackageMappings[packageConstant]; exists {
		return mapping.Debian
	}
	return packageConstant // Fallback to original name if not mapped
}
