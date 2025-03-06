package playwrightmatchers_test

import (
	"github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("HaveAttribute", func() {

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
		matcher := playwrightmatchers.HaveAttribute("class", "expected-class")
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
		matcher := playwrightmatchers.HaveAttribute("class", "expected-class")
		success, err := matcher.Match(locator)
		Expect(success).To(BeFalse())
		Expect(err).To(BeNil())

		failureMessage := matcher.FailureMessage(locator)
		Expect(failureMessage).To(ContainSubstring("to have attribute class=expected-class"))

		normalizedMessage := normalizeHTML(failureMessage)
		Expect(normalizedMessage).To(ContainSubstring(`<div id="test" class="wrong-class">Test Element</div>`))
	})

	It("should return error for invalid type", func() {
		matcher := playwrightmatchers.HaveAttribute("class", "test")
		success, err := matcher.Match("invalid")
		Expect(success).To(BeFalse())
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("expected playwright.Locator, got string"))
	})
})
