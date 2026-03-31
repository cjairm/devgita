package embedded

// ExtractFunc is a function type for extracting embedded configs
// This allows dependency injection for testing
type ExtractFunc func(destDir string) error

// DefaultExtractor will be set by main package
var DefaultExtractor ExtractFunc
