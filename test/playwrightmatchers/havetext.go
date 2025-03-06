package playwrightmatchers

import (
	"fmt"
	"time"

	"github.com/onsi/gomega/types"
	"github.com/playwright-community/playwright-go"
	"github.com/yosssi/gohtml"
)

type haveTextMatcher struct {
	text    string
	timeout time.Duration
}

func (h haveTextMatcher) Match(actual any) (success bool, err error) {
	if page, err := toPage(actual); err == nil {
		err = expect.Locator(page.GetByText(h.text)).ToBeVisible(playwright.LocatorAssertionsToBeVisibleOptions{
			Timeout: playwright.Float(float64(h.timeout.Milliseconds())),
		})
		return err == nil, nil
	}

	locator, err := toLocator(actual)
	if err != nil {
		return false, fmt.Errorf("expected playwright.Page or playwright.Locator, got %T", actual)
	}

	err = expect.Locator(locator).ToHaveText(h.text, playwright.LocatorAssertionsToHaveTextOptions{
		Timeout: playwright.Float(float64(h.timeout.Milliseconds())),
	})

	if err != nil {
		page, err := locator.Page()
		if err != nil {
			return false, fmt.Errorf("cannot retrieve page: %w", err)
		}
		content, err := page.Content()
		if err != nil {
			return false, fmt.Errorf("cannot retrieve page content: %w", err)
		}

		fmt.Println(gohtml.Format(content))
	}
	return err == nil, nil
}

func (h haveTextMatcher) FailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain text \n\t%#v\nwithin %v", actual, h.text, h.timeout)
}

func (h haveTextMatcher) NegatedFailureMessage(actual any) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to contain text \n\t%#v\nwithin %v", actual, h.text, h.timeout)
}

func HaveText(text string, timeout ...time.Duration) types.GomegaMatcher {
	return &haveTextMatcher{
		text:    text,
		timeout: getTimeout(timeout...),
	}
}
