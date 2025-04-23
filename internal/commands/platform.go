package commands

import "runtime"

type CustomizablePlatform interface {
	IsLinux() bool
	IsMac() bool
}

type Platform struct{}

func NewPlatform() *Platform {
	return &Platform{}
}

func (Platform) IsLinux() bool {
	return runtime.GOOS == "linux"
}

func (Platform) IsMac() bool {
	return runtime.GOOS == "darwin"
}
