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

func (h *haveCount) Match(actual any) (success bool, err error) {
	locator, err := toLocator(actual)
	if err != nil {
		return false, err
	}

	err = expect.Locator(locator).ToHaveCount(h.expected, playwright.LocatorAssertionsToHaveCountOptions{
		Timeout: playwright.Float(float64(h.timeout.Milliseconds())),
	})
	locators, _ := locator.All()
	h.actual = len(locators)

	return err == nil, nil
}

func (h *haveCount) FailureMessage(actual any) (message string) {
	content, err := PrettyPrintHTML(actual)
	if err == nil {
		return fmt.Sprintf(`Expected locator to have count of %d
Actual locator has count %d
within %v
Page content:
%s`, h.expected, h.actual, h.timeout, content)
	}
	return fmt.Sprintf(`Expected locator to have count of %d
Actual locator has count %d
within %v`, h.expected, h.actual, h.timeout)
}

func (h *haveCount) NegatedFailureMessage(actual any) (message string) {
	content, err := PrettyPrintHTML(actual)
	if err == nil {
		return fmt.Sprintf(`Expected locator to not have count %d
Actual locator has count of %d
within %v
Page content:
%s`, h.expected, h.actual, h.timeout, content)
	}
	return fmt.Sprintf(`Expected locator to not have count %d
Actual locator has count of %d
within %v`, h.expected, h.actual, h.timeout)
}

func HaveCount(expected int, timeout ...time.Duration) types.GomegaMatcher {
	return &haveCount{
		expected: expected,
		timeout:  getTimeout(timeout...),
	}
}
