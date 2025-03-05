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

func (b beVisibleMatcher) Match(actual interface{}) (success bool, err error) {
	locator, ok := actual.(playwright.Locator)
	if !ok {
		return false, fmt.Errorf("BeVisible matcher expects a playwright.Locator")
	}

	// First check if the element exists
	exists, err := locator.Count()
	if err != nil {
		return false, err
	}
	if exists == 0 {
		return false, nil
	}

	// Use configured timeout for visibility check
	err = expect.Locator(locator).ToBeVisible(playwright.LocatorAssertionsToBeVisibleOptions{
		Timeout: playwright.Float(float64(b.timeout.Milliseconds())),
	})
	return err == nil, nil
}

func (b beVisibleMatcher) FailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected element to be visible within %v", b.timeout)
}

func (b beVisibleMatcher) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf("Expected element to not be visible within %v", b.timeout)
}

// BeVisible creates a matcher that checks if an element is visible.
// The timeout parameter is optional and defaults to 500ms.
func BeVisible(timeout ...time.Duration) types.GomegaMatcher {
	t := GetDefaultTimeout()
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return &beVisibleMatcher{timeout: t}
}
