package playwrightmatchers

import (
	"fmt"

	"github.com/onsi/gomega/types"
	"github.com/playwright-community/playwright-go"
)

type haveAttributeMatcher struct {
	name  string
	value string
}

func (h haveAttributeMatcher) Match(actual any) (success bool, err error) {
	locator, ok := actual.(playwright.Locator)
	if !ok {
		return false, fmt.Errorf("HaveAttribute matcher expects a playwright.Locator")
	}

	// First check if the element exists
	exists, err := locator.Count()
	if err != nil {
		return false, err
	}
	if exists == 0 {
		return false, nil
	}

	value, err := locator.GetAttribute(h.name)
	if err != nil {
		return false, err
	}

	return value == h.value, nil
}

func (h haveAttributeMatcher) FailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto have attribute %s=%s", actual, h.name, h.value)
}

func (h haveAttributeMatcher) NegatedFailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to have attribute %s=%s", actual, h.name, h.value)
}

func HaveAttribute(name, value string) types.GomegaMatcher {
	return &haveAttributeMatcher{
		name:  name,
		value: value,
	}
}
