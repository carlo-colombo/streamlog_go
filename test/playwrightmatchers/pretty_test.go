package playwrightmatchers_test

import (
	"strings"

	"github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

// normalizeHTML removes extra whitespace, newlines, and tabs from HTML content
func normalizeHTML(html string) string {
	// Replace newlines and tabs with spaces
	html = strings.ReplaceAll(html, "\n", " ")
	html = strings.ReplaceAll(html, "\t", " ")
	// Replace multiple spaces with single space
	html = strings.Join(strings.Fields(html), " ")
	// Trim spaces around tags
	html = strings.ReplaceAll(html, "> ", ">")
	html = strings.ReplaceAll(html, " <", "<")
	// Trim spaces around text content
	html = strings.ReplaceAll(html, ">  ", "> ")
	html = strings.ReplaceAll(html, "  <", " <")
	// Trim leading/trailing whitespace
	return strings.TrimSpace(html)
}

var _ = ginkgo.Describe("PrettyPrintHTML", func() {
	var pw *playwright.Playwright
	var browser playwright.Browser
	var page playwright.Page

	ginkgo.BeforeEach(func() {
		var err error
		pw, err = playwright.Run()
		gomega.Expect(err).To(gomega.BeNil())
		browser, err = pw.Chromium.Launch()
		gomega.Expect(err).To(gomega.BeNil())
		page, err = browser.NewPage()
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.AfterEach(func() {
		page.Close()
		browser.Close()
		pw.Stop()
	})

	ginkgo.It("should format HTML content from a Page", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div id="test">Hello World</div>
				</body>
			</html>
		`)
		gomega.Expect(err).To(gomega.BeNil())

		content, err := playwrightmatchers.PrettyPrintHTML(page)
		gomega.Expect(err).To(gomega.BeNil())
		normalizedContent := normalizeHTML(content)
		gomega.Expect(normalizedContent).To(gomega.ContainSubstring(`<div id="test">Hello World</div>`))
	})

	ginkgo.It("should format HTML content from a Locator", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div id="test">Hello World</div>
				</body>
			</html>
		`)
		gomega.Expect(err).To(gomega.BeNil())

		locator := page.Locator("#test")
		content, err := playwrightmatchers.PrettyPrintHTML(locator)
		gomega.Expect(err).To(gomega.BeNil())
		normalizedContent := normalizeHTML(content)
		gomega.Expect(normalizedContent).To(gomega.ContainSubstring(`<div id="test">Hello World</div>`))
	})

	ginkgo.It("should return error for invalid type", func() {
		_, err := playwrightmatchers.PrettyPrintHTML("invalid")
		gomega.Expect(err).To(gomega.HaveOccurred())
		gomega.Expect(err.Error()).To(gomega.ContainSubstring("unsupported type for pretty printing"))
	})
})
