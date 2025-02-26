package playwrightmatchers

import (
	"fmt"
	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"
)

var expect = playwright.NewPlaywrightAssertions()

type haveCount struct {
	expected int
	actual   int
}

func (h *haveCount) Match(actual interface{}) (success bool, err error) {
	actualLocator, ok := actual.(playwright.Locator)

	if !ok {
		return false, fmt.Errorf("HaveCount matcher expects a playwright.Locator")
	}

	err = expect.Locator(actualLocator).ToHaveCount(h.expected)
	locators, _ := actualLocator.All()
	h.actual = len(locators)

	return err == nil, nil
}

func (h *haveCount) FailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf(`Expected locator to have count of %d
Actual locator has count %d`, h.expected, h.actual)
}

func (h *haveCount) NegatedFailureMessage(_ interface{}) (message string) {
	return fmt.Sprintf(`Expected locator to not have count %d
Actual locator has count of %d`, h.expected, h.actual)
}

func HaveCount(expected int) types.GomegaMatcher {
	return &haveCount{expected: expected}
}

func HaveText(text string) types.GomegaMatcher {
	return &haveTextMatcher{text}
}

func HaveSelector(selector string) types.GomegaMatcher {
	return &haveSelectorMatcher{selector: selector}
}

type haveSelectorMatcher struct {
	selector string
}
type haveTextMatcher struct {
	text string
}

func (h haveTextMatcher) Match(actual interface{}) (success bool, err error) {
	if page, ok := actual.(playwright.Page); ok {
		err = expect.Locator(page.GetByText(h.text)).ToBeVisible()
		return err == nil, nil
	}

	if locator, ok := actual.(playwright.Locator); ok {
		err = expect.Locator(locator).ToHaveText(h.text)
		return err == nil, nil
	}
	return false, errors.New("HaveTextMatcher expects a playwright.Page or playwright.Locator")
}

func (h haveTextMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain text \n\t%#v", actual, h.text)
}

func (h haveTextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\n not to contain text \n\t%#v", actual, h.text)
}

func (h haveSelectorMatcher) Match(actual interface{}) (success bool, err error) {
	page, ok := actual.(playwright.Page)

	if !ok {
		return false, fmt.Errorf("HaveSelector matcher expects a playwright.Page")
	}

	err = expect.Locator(page.Locator(h.selector)).ToBeVisible()

	return err == nil, err
}

func (h haveSelectorMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain selector \n\t%#v", actual, h.selector)
}

func (h haveSelectorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\n not to contain selector \n\t%#v", actual, h.selector)
}
