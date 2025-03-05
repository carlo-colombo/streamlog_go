package playwrightmatchers_test

import (
	"time"

	"github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("BeVisible", func() {
	var page playwright.Page

	BeforeEach(func() {
		var err error
		page, err = browser.NewPage()
		Expect(err).Should(BeNil())
	})

	AfterEach(func() {
		Expect(page.Close()).Should(BeNil())
	})

	It("should match when element is visible", func() {
		err := page.SetContent(`
			<div id="visible" style="display: block;">Visible Element</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator("#visible")
		Expect(locator).Should(playwrightmatchers.BeVisible())
	})

	It("should not match when element is hidden", func() {
		err := page.SetContent(`
			<div id="hidden" style="display: none;">Hidden Element</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator("#hidden")
		Expect(locator).ShouldNot(playwrightmatchers.BeVisible(1 * time.Millisecond))
	})

	It("should not match when element doesn't exist", func() {
		locator := page.Locator("#non-existent")
		Expect(locator).ShouldNot(playwrightmatchers.BeVisible(1 * time.Millisecond))
	})

	It("should match when element becomes visible", func() {
		err := page.SetContent(`
			<div id="dynamic" style="display: none;">Dynamic Element</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator("#dynamic")
		Expect(locator).ShouldNot(playwrightmatchers.BeVisible(1 * time.Millisecond))

		_, err = page.Evaluate(`() => {
			document.getElementById('dynamic').style.display = 'block';
		}`)
		Expect(err).Should(BeNil())

		Expect(locator).Should(playwrightmatchers.BeVisible())
	})

	It("should fail with appropriate message when element is not visible", func() {
		err := page.SetContent(`
			<div id="hidden" style="display: none;">Hidden Element</div>
		`)
		Expect(err).Should(BeNil())

		locator := page.Locator("#hidden")

		success, err := playwrightmatchers.BeVisible(1 * time.Millisecond).Match(locator)
		Expect(err).Should(BeNil())
		Expect(success).Should(BeFalse())
	})
})
