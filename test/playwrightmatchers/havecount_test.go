package playwrightmatchers_test

import (
	"github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("HaveCount", func() {
	var page playwright.Page

	BeforeEach(func() {
		var err error
		page, err = browser.NewPage()
		Expect(err).Should(BeNil())
	})

	AfterEach(func() {
		Expect(page.Close()).Should(BeNil())
	})

	It("should match when element count matches expected", func() {
		err := page.SetContent(`
			<div class="item">Item 1</div>
			<div class="item">Item 2</div>
			<div class="item">Item 3</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator(".item")
		Expect(locator).Should(playwrightmatchers.HaveCount(3))
	})

	It("should not match when element count is different", func() {
		err := page.SetContent(`
			<div class="item">Item 1</div>
			<div class="item">Item 2</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator(".item")
		Expect(locator).ShouldNot(playwrightmatchers.HaveCount(3))
	})

	It("should not match when no elements exist", func() {
		locator := page.Locator(".non-existent")
		Expect(locator).ShouldNot(playwrightmatchers.HaveCount(1))
	})

	It("should fail with appropriate message when count is wrong", func() {
		err := page.SetContent(`
			<div class="item">Item 1</div>
			<div class="item">Item 2</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator(".item")
		success, err := playwrightmatchers.HaveCount(3).Match(locator)
		Expect(err).Should(BeNil())
		Expect(success).Should(BeFalse())
	})
})
