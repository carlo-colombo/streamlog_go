package playwrightmatchers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("HaveAttribute", func() {
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

	It("should match when element has the specified attribute", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div id="test" class="expected-class">Test Element</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator("#test")
		matcher := HaveAttribute("class", "expected-class")
		success, err := matcher.Match(locator)
		Expect(success).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("should not match when attribute value is different", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div id="test" class="wrong-class">Test Element</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator("#test")
		matcher := HaveAttribute("class", "expected-class")
		success, err := matcher.Match(locator)
		Expect(success).To(BeFalse())
		Expect(err).To(BeNil())

		failureMessage := matcher.FailureMessage(locator)
		Expect(failureMessage).To(ContainSubstring("to have attribute class=expected-class"))

		normalizedMessage := normalizeHTML(failureMessage)
		Expect(normalizedMessage).To(ContainSubstring(`<div id="test" class="wrong-class">Test Element</div>`))
	})

	It("should return error for invalid type", func() {
		matcher := HaveAttribute("class", "test")
		success, err := matcher.Match("invalid")
		Expect(success).To(BeFalse())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected playwright.Locator, got string"))
	})
})
