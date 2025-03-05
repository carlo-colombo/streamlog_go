package playwrightmatchers

import "time"

var defaultTimeout = 500 * time.Millisecond

// SetDefaultTimeout sets the default timeout for all matchers
func SetDefaultTimeout(timeout time.Duration) {
	defaultTimeout = timeout
}

// GetDefaultTimeout returns the current default timeout
func GetDefaultTimeout() time.Duration {
	return defaultTimeout
}
