package playwrightmatchers

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("HaveCount", func() {
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

	It("should include HTML content in failure message", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div class="test">First Element</div>
					<div class="test">Second Element</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator(".test")
		matcher := HaveCount(3)
		success, err := matcher.Match(locator)
		Expect(success).To(BeFalse())
		Expect(err).To(BeNil())

		failureMessage := matcher.FailureMessage(locator)
		Expect(failureMessage).To(ContainSubstring("Expected locator to have count of 3"))
		Expect(failureMessage).To(ContainSubstring("Actual locator has count 2"))

		normalizedMessage := normalizeHTML(failureMessage)
		Expect(normalizedMessage).To(ContainSubstring(`<div class="test">First Element</div>`))
		Expect(normalizedMessage).To(ContainSubstring(`<div class="test">Second Element</div>`))
	})

	It("should match when element count matches expected", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div class="test">First Element</div>
					<div class="test">Second Element</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator(".test")
		matcher := HaveCount(2)
		success, err := matcher.Match(locator)
		Expect(success).To(BeTrue())
		Expect(err).To(BeNil())
	})

	It("should not match when element count is different", func() {
		err := page.SetContent(`
			<html>
				<body>
					<div class="test">First Element</div>
					<div class="test">Second Element</div>
				</body>
			</html>
		`)
		Expect(err).To(BeNil())

		locator := page.Locator(".test")
		matcher := HaveCount(3)
		success, err := matcher.Match(locator)
		Expect(success).To(BeFalse())
		Expect(err).To(BeNil())

		failureMessage := matcher.FailureMessage(locator)
		Expect(failureMessage).To(ContainSubstring("Expected locator to have count of 3"))
		Expect(failureMessage).To(ContainSubstring("Actual locator has count 2"))

		normalizedMessage := normalizeHTML(failureMessage)
		Expect(normalizedMessage).To(ContainSubstring(`<div class="test">First Element</div>`))
		Expect(normalizedMessage).To(ContainSubstring(`<div class="test">Second Element</div>`))
	})

	It("should return error for invalid type", func() {
		matcher := HaveCount(1)
		success, err := matcher.Match("invalid")
		Expect(success).To(BeFalse())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected playwright.Locator, got string"))
	})
})
