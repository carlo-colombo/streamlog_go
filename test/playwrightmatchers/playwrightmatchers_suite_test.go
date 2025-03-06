package playwrightmatchers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var pw *playwright.Playwright
var browser playwright.Browser
var page playwright.Page

func TestPlaywrightmatchers(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var err error
		pw, err = playwright.Run()
		Expect(err).Should(BeNil())

		browser, err = pw.Chromium.Launch()
		Expect(err).Should(BeNil())
	})

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

	AfterSuite(func() {
		Expect(browser.Close()).Should(BeNil())
		Expect(pw.Stop()).Should(BeNil())
	})

	RunSpecs(t, "BeVisible Matcher Suite")
}
