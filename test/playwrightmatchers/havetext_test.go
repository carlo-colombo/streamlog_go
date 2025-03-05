package playwrightmatchers_test

import (
	"github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("HaveText", func() {
	var page playwright.Page

	BeforeEach(func() {
		var err error
		page, err = browser.NewPage()
		Expect(err).Should(BeNil())
	})

	AfterEach(func() {
		Expect(page.Close()).Should(BeNil())
	})

	It("should match when page contains the text", func() {
		err := page.SetContent(`
			<div>Hello World</div>
		`)
		Expect(err).Should(BeNil())

		Expect(page).Should(playwrightmatchers.HaveText("Hello World"))
	})

	It("should match when locator contains the text", func() {
		err := page.SetContent(`
			<div id="test">Hello World</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator("#test")
		Expect(locator).Should(playwrightmatchers.HaveText("Hello World"))
	})

	It("should not match when text is not present", func() {
		err := page.SetContent(`
			<div>Different Text</div>
		`)
		Expect(err).Should(BeNil())

		Expect(page).ShouldNot(playwrightmatchers.HaveText("Hello World"))
	})

	It("should not match when element doesn't exist", func() {
		locator := page.Locator("#non-existent")
		Expect(locator).ShouldNot(playwrightmatchers.HaveText("Hello World"))
	})

	It("should fail with appropriate message when text is not found", func() {
		err := page.SetContent(`
			<div>Different Text</div>
		`)
		Expect(err).Should(BeNil())

		success, err := playwrightmatchers.HaveText("Hello World").Match(page)
		Expect(err).Should(BeNil())
		Expect(success).Should(BeFalse())
	})
})
