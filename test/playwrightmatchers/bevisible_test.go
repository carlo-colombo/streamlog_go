package playwrightmatchers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("BeVisible", func() {
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

	It("should match when element is visible", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div id="test">Visible Element</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator("#test")
		matcher := BeVisible()
		success, err := matcher.Match(locator)
		Expect(success).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("should not match when element is hidden", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div id="test" style="display: none">Hidden Element</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator("#test")
		matcher := BeVisible()
		success, err := matcher.Match(locator)
		Expect(success).To(BeFalse())
		Expect(err).To(BeNil())

		failureMessage := matcher.FailureMessage(locator)
		Expect(failureMessage).To(ContainSubstring("Expected element to be visible"))

		normalizedMessage := normalizeHTML(failureMessage)
		Expect(normalizedMessage).To(ContainSubstring(`<div id="test" style="display: none">Hidden Element</div>`))
	})

	It("should return error for invalid type", func() {
		matcher := BeVisible()
		success, err := matcher.Match("invalid")
		Expect(success).To(BeFalse())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected playwright.Locator, got string"))
	})
})
