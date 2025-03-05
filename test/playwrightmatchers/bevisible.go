package playwrightmatchers

import (
	"fmt"
	"time"

	"github.com/onsi/gomega/types"
	"github.com/playwright-community/playwright-go"
)

type beVisibleMatcher struct {
	timeout time.Duration
}

func (b beVisibleMatcher) Match(actual any) (success bool, err error) {
	locator, err := toLocator(actual)
	if err != nil {
		return false, err
	}

	err = expect.Locator(locator).ToBeVisible(playwright.LocatorAssertionsToBeVisibleOptions{
		Timeout: playwright.Float(float64(b.timeout.Milliseconds())),
	})
	return err == nil, nil
}

func (b beVisibleMatcher) FailureMessage(_ any) (message string) {
	return fmt.Sprintf("Expected element to be visible within %v", b.timeout)
}

func (b beVisibleMatcher) NegatedFailureMessage(_ any) (message string) {
	return fmt.Sprintf("Expected element to not be visible within %v", b.timeout)
}

// BeVisible creates a matcher that checks if an element is visible.
// The timeout parameter is optional and defaults to 500ms.
func BeVisible(timeout ...time.Duration) types.GomegaMatcher {
	return &beVisibleMatcher{timeout: getTimeout(timeout...)}
}
