package playwrightmatchers

import (
	"fmt"

	"github.com/playwright-community/playwright-go"
	"github.com/yosssi/gohtml"
)

// PrettyPrintHTML retrieves and formats the HTML content of a page or locator
func PrettyPrintHTML(actual any) (string, error) {
	var page playwright.Page
	var err error

	switch v := actual.(type) {
	case playwright.Page:
		page = v
	case playwright.Locator:
		page, err = v.Page()
		if err != nil {
			return "", fmt.Errorf("cannot retrieve page: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported type for pretty printing: %T", actual)
	}

	content, err := page.Content()
	if err != nil {
		return "", fmt.Errorf("cannot retrieve page content: %w", err)
	}

	return gohtml.Format(content), nil
}
