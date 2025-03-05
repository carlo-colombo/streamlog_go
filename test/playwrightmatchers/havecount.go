package playwrightmatchers

import (
	"fmt"
	"time"

	"github.com/onsi/gomega/types"
	"github.com/playwright-community/playwright-go"
)

type haveCount struct {
	expected int
	actual   int
	timeout  time.Duration
}

func (h *haveCount) Match(actual interface{}) (success bool, err error) {
	actualLocator, ok := actual.(playwright.Locator)

	if !ok {
		return false, fmt.Errorf("HaveCount matcher expects a playwright.Locator")
	}

	err = expect.Locator(actualLocator).ToHaveCount(h.expected, playwright.LocatorAssertionsToHaveCountOptions{
		Timeout: playwright.Float(float64(h.timeout.Milliseconds())),
	})
	locators, _ := actualLocator.All()
	h.actual = len(locators)

	return err == nil, nil
}

func (h *haveCount) FailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf(`Expected locator to have count of %d
Actual locator has count %d
within %v`, h.expected, h.actual, h.timeout)
}

func (h *haveCount) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf(`Expected locator to not have count %d
Actual locator has count of %d
within %v`, h.expected, h.actual, h.timeout)
}

func HaveCount(expected int, timeout ...time.Duration) types.GomegaMatcher {
	t := GetDefaultTimeout()
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return &haveCount{expected: expected, timeout: t}
}
