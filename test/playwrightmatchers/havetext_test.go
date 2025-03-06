package playwrightmatchers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("HaveText", func() {
	var pw *playwright.Playwright
	var browser playwright.Browser
	var page playwright.Page

	BeforeEach(func() {
		var err error
		pw, err = playwright.Run()
		Expect(err).To(BeNil())
		browser, err = pw.Chromium.Launch()
		Expect(err).To(BeNil())
		page, err = browser.NewPage()
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		page.Close()
		browser.Close()
		pw.Stop()
	})

	It("should match when element contains the expected text", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div id="test">Expected Text</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator("#test")
		matcher := HaveText("Expected Text")
		success, err := matcher.Match(locator)
		Expect(success).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("should match when page contains the expected text", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div>Expected Text</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		matcher := HaveText("Expected Text")
		success, err := matcher.Match(page)
		Expect(success).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("should not match when text is different", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div id="test">Wrong Text</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator("#test")
		matcher := HaveText("Expected Text")
		success, err := matcher.Match(locator)
		Expect(success).To(BeFalse())
		Expect(err).To(BeNil())

		failureMessage := matcher.FailureMessage(locator)
		Expect(failureMessage).To(ContainSubstring("to contain text"))
		Expect(failureMessage).To(ContainSubstring("Expected Text"))

		normalizedMessage := normalizeHTML(failureMessage)
		Expect(normalizedMessage).To(ContainSubstring(`<div id="test">Wrong Text</div>`))
	})

	It("should return error for invalid type", func() {
		matcher := HaveText("test")
		success, err := matcher.Match("invalid")
		Expect(success).To(BeFalse())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected playwright.Page or playwright.Locator"))
	})
})
