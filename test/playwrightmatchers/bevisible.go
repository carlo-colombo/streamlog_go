package playwrightmatchers

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"github.com/playwright-community/playwright-go"
)

type beVisibleMatcher struct {
	timeout float64
}

func (b beVisibleMatcher) Match(actual any) (success bool, err error) {
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
		Timeout: playwright.Float(b.timeout),
	})
	return err == nil, nil
}

func (b beVisibleMatcher) FailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto be visible", actual)
}

func (b beVisibleMatcher) NegatedFailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to be visible", actual)
}

// BeVisible creates a matcher that checks if an element is visible.
// The timeout parameter is optional and defaults to 500ms.
func BeVisible(timeout ...float64) types.GomegaMatcher {
	t := 500.0 // default timeout in milliseconds
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return &beVisibleMatcher{
		timeout: t,
	}
}
