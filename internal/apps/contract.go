package apps

// AppKind classifies what kind of application an app is.
type AppKind int

const (
	KindUnknown  AppKind = iota
	KindTerminal         // CLI tools, terminal emulators, shell utilities
	KindDesktop          // GUI desktop applications
	KindLanguage         // Programming language runtimes and managers
	KindDatabase         // Database systems
	KindFont             // Font packages (satisfies FontInstaller, not App)
	KindMeta             // Devgita itself
)

// App is the contract every app module must satisfy.
// Fonts is the only exception — it satisfies FontInstaller instead.
type App interface {
	Name() string
	Kind() AppKind

	Install() error
	ForceInstall() error
	SoftInstall() error

	ForceConfigure() error
	SoftConfigure() error

	Uninstall() error
	Update() error

	ExecuteCommand(args ...string) error
}

// SelectiveConfigurer is an optional interface for apps whose configuration
// includes discrete, separately-refreshable subtrees — the shared
// skills/commands/agents trees for the AI coders (claude, opencode). It backs
// `dg configure <app> --force --only=...`, letting a user overwrite just those
// subtrees without disturbing general config they may have edited (settings,
// themes, generated files). Apps that don't implement it reject --only.
type SelectiveConfigurer interface {
	// ConfigurableParts lists the part names accepted by --only.
	ConfigurableParts() []string
	// ForceConfigureParts overwrites only the named parts, leaving all other
	// configuration in place.
	ForceConfigureParts(parts []string) error
}

// FontInstaller is the contract for the Fonts module, which installs named fonts
// rather than a single application.
type FontInstaller interface {
	Name() string
	Kind() AppKind
	Available() []string
	SoftInstallAll()
	InstallFont(name string) error
	ForceInstallFont(name string) error
	SoftInstallFont(name string) error
	UninstallFont(name string) error
}
