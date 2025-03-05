package playwrightmatchers

import (
	"fmt"
	"time"

	"github.com/onsi/gomega/types"
	"github.com/pkg/errors"
	"github.com/playwright-community/playwright-go"
)

type haveTextMatcher struct {
	text    string
	timeout time.Duration
}

func (h haveTextMatcher) Match(actual interface{}) (success bool, err error) {
	if page, ok := actual.(playwright.Page); ok {
		err = expect.Locator(page.GetByText(h.text)).ToBeVisible(playwright.LocatorAssertionsToBeVisibleOptions{
			Timeout: playwright.Float(float64(h.timeout.Milliseconds())),
		})
		return err == nil, nil
	}

	if locator, ok := actual.(playwright.Locator); ok {
		err = expect.Locator(locator).ToHaveText(h.text, playwright.LocatorAssertionsToHaveTextOptions{
			Timeout: playwright.Float(float64(h.timeout.Milliseconds())),
		})
		return err == nil, nil
	}
	return false, errors.New("HaveTextMatcher expects a playwright.Page or playwright.Locator")
}

func (h haveTextMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain text \n\t%#v\nwithin %v", actual, h.text, h.timeout)
}

func (h haveTextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to contain text \n\t%#v\nwithin %v", actual, h.text, h.timeout)
}

func HaveText(text string, timeout ...time.Duration) types.GomegaMatcher {
	t := GetDefaultTimeout()
	if len(timeout) > 0 {
		t = timeout[0]
	}
	return &haveTextMatcher{text: text, timeout: t}
}
