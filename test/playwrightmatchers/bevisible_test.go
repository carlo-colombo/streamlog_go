package playwrightmatchers_test

import (
	"github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BeVisible", func() {

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
		matcher := playwrightmatchers.BeVisible()
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
		matcher := playwrightmatchers.BeVisible()
		success, err := matcher.Match(locator)
		Expect(success).To(BeFalse())
		Expect(err).To(BeNil())

		failureMessage := matcher.FailureMessage(locator)
		Expect(failureMessage).To(ContainSubstring("Expected element to be visible"))

		normalizedMessage := normalizeHTML(failureMessage)
		Expect(normalizedMessage).To(ContainSubstring(`<div id="test" style="display: none">Hidden Element</div>`))
	})

	It("should return error for invalid type", func() {
		matcher := playwrightmatchers.BeVisible()
		success, err := matcher.Match("invalid")
		Expect(success).To(BeFalse())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected playwright.Locator, got string"))
	})
})
