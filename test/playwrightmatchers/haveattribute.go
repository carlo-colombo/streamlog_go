package playwrightmatchers

import (
	"fmt"
	"time"

	"github.com/onsi/gomega/types"
	"github.com/playwright-community/playwright-go"
)

type haveAttributeMatcher struct {
	name    string
	value   string
	timeout time.Duration
}

func (h haveAttributeMatcher) Match(actual any) (success bool, err error) {
	locator, err := toLocator(actual)
	if err != nil {
		return false, err
	}

	err = expect.Locator(locator).ToHaveAttribute(h.name, h.value, playwright.LocatorAssertionsToHaveAttributeOptions{
		Timeout: playwright.Float(float64(h.timeout.Milliseconds())),
	})
	return err == nil, nil
}

func (h haveAttributeMatcher) FailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto have attribute %s=%s within %v", actual, h.name, h.value, h.timeout)
}

func (h haveAttributeMatcher) NegatedFailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to have attribute %s=%s within %v", actual, h.name, h.value, h.timeout)
}

func HaveAttribute(name, value string, timeout ...time.Duration) types.GomegaMatcher {
	return &haveAttributeMatcher{
		name:    name,
		value:   value,
		timeout: getTimeout(timeout...),
	}
}
