package playwrightmatchers_test

import (
	"github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("HaveAttribute", func() {
	var page playwright.Page

	BeforeEach(func() {
		var err error
		page, err = browser.NewPage()
		Expect(err).Should(BeNil())
	})

	AfterEach(func() {
		Expect(page.Close()).Should(BeNil())
	})

	It("should match when element has the specified attribute", func() {
		err := page.SetContent(`
			<input type="text" name="test" placeholder="Test input">
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator("input")
		Expect(locator).Should(playwrightmatchers.HaveAttribute("type", "text"))
		Expect(locator).Should(playwrightmatchers.HaveAttribute("name", "test"))
		Expect(locator).Should(playwrightmatchers.HaveAttribute("placeholder", "Test input"))
	})

	It("should not match when element doesn't have the attribute", func() {
		err := page.SetContent(`
			<div>No attributes</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator("div")
		Expect(locator).ShouldNot(playwrightmatchers.HaveAttribute("type", "text"))
	})

	It("should not match when attribute value is different", func() {
		err := page.SetContent(`
			<input type="password" name="test">
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator("input")
		Expect(locator).ShouldNot(playwrightmatchers.HaveAttribute("type", "text"))
	})

	It("should not match when element doesn't exist", func() {
		locator := page.Locator("#non-existent")
		Expect(locator).ShouldNot(playwrightmatchers.HaveAttribute("type", "text"))
	})
})
