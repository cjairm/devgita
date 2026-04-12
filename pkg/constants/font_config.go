package constants

import "fmt"

// NerdFontConfig represents the configuration for a Nerd Font across platforms
type NerdFontConfig struct {
	DisplayName string // Human-readable font name (e.g., "Hack Nerd Font")
	PackageName string // macOS Homebrew cask name (e.g., "font-hack-nerd-font")
	ArchiveName string // GitHub nerd-fonts release archive name without extension (e.g., "Hack")
	InstallName string // Font detection string for fc-list lookup (e.g., "Hack")
}

// GetFontConfigs returns the list of supported Nerd Fonts
func GetFontConfigs() []NerdFontConfig {
	return []NerdFontConfig{
		{
			DisplayName: "Hack Nerd Font",
			PackageName: "font-hack-nerd-font",
			ArchiveName: "Hack",
			InstallName: "Hack Nerd Font",
		},
		{
			DisplayName: "Meslo LG Nerd Font",
			PackageName: "font-meslo-lg-nerd-font",
			ArchiveName: "Meslo",
			InstallName: "MesloLGS NF",
		},
		{
			DisplayName: "CaskaydiaMono Nerd Font",
			PackageName: "font-caskaydia-mono-nerd-font",
			ArchiveName: "CascadiaMono",
			InstallName: "CaskaydiaMono Nerd Font",
		},
		{
			DisplayName: "FiraMono Nerd Font",
			PackageName: "font-fira-mono-nerd-font",
			ArchiveName: "FiraMono",
			InstallName: "FiraMono Nerd Font",
		},
		{
			DisplayName: "JetBrainsMono Nerd Font",
			PackageName: "font-jetbrains-mono-nerd-font",
			ArchiveName: "JetBrainsMono",
			InstallName: "JetBrainsMono Nerd Font",
		},
	}
}

// GetNerdFontURL constructs the GitHub release download URL for a Nerd Font archive
func GetNerdFontURL(archiveName string) string {
	return fmt.Sprintf(
		"https://github.com/ryanoasis/nerd-fonts/releases/latest/download/%s.tar.xz",
		archiveName,
	)
}

// GetFontConfigByPackageName looks up a NerdFontConfig by its Homebrew package name
// Returns nil if no matching config is found
func GetFontConfigByPackageName(packageName string) *NerdFontConfig {
	for _, fc := range GetFontConfigs() {
		if fc.PackageName == packageName {
			return &fc
		}
	}
	return nil
}
