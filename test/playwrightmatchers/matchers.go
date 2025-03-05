package playwrightmatchers

import (
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"
)

const defaultTimeout = 500 * time.Millisecond

// GetDefaultTimeout returns the default timeout for matchers
func GetDefaultTimeout() time.Duration {
	return defaultTimeout
}

// toLocator safely converts an any to a playwright.Locator
func toLocator(actual any) (playwright.Locator, error) {
	locator, ok := actual.(playwright.Locator)
	if !ok {
		return nil, fmt.Errorf("expected playwright.Locator, got %T", actual)
	}
	return locator, nil
}

// toPage safely converts an any to a playwright.Page
func toPage(actual any) (playwright.Page, error) {
	page, ok := actual.(playwright.Page)
	if !ok {
		return nil, fmt.Errorf("expected playwright.Page, got %T", actual)
	}
	return page, nil
}

// getTimeout returns the timeout from the variadic parameter or the default timeout
func getTimeout(timeout ...time.Duration) time.Duration {
	if len(timeout) > 0 {
		return timeout[0]
	}
	return GetDefaultTimeout()
}

var expect = playwright.NewPlaywrightAssertions()
