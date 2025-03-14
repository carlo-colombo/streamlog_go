package integration_test

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
)

var _ = Describe("e2e tests", func() {
	BeforeEach(func() { PauseOutputInterception() })

	AfterEach(func() { ResumeOutputInterception() })

	Describe("the page shown as the /", func() {
		It("shows logs", func() {
			By("testing if it is available")
			resp, err := http.Get(targetUrl)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(resp).To(SatisfyAll(
				HaveHTTPStatus(http.StatusOK),
				HaveHTTPHeaderWithValue("Content-Type", ContainSubstring("text/html")),
			))

			By("starting a browser")
			pw, err := playwright.Run()
			Expect(err).ShouldNot(HaveOccurred())

			browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
				Headless: playwright.Bool(false),
			})
			Expect(err).ShouldNot(HaveOccurred())

			page, err := browser.NewPage()
			Expect(err).ShouldNot(HaveOccurred())

			By("opening the web browser")
			_, err = page.Goto(targetUrl)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(page.Locator("h1")).To(BeVisible())
			Expect(page).To(HaveText("Streamlog"))

			lines := []string{
				"line from stdin",
				"and another",
				"bonus line from stdin",
				"a new message that should show stdin",
				"another new message that should not show",
			}

			_, _ = fmt.Fprintln(stdinWriter, lines[0])
			_, _ = fmt.Fprintln(stdinWriter, lines[1])

			Expect(page).To(HaveText(lines[0]))
			Expect(page).To(HaveText(lines[1]))

			Expect(page.Locator("table tr")).To(HaveCount(2))

			By("sending an additional line to stdin and prepending to the content")
			_, _ = fmt.Fprintln(stdinWriter, lines[2])
			Expect(page.Locator("table tr:first-child td.message")).To(HaveText(lines[2]))

			By("setting a filter")
			filterInput := page.Locator("input[placeholder='Filter logs...']")
			Expect(filterInput).To(BeVisible())

			err = filterInput.Fill("stdin")
			Expect(err).ShouldNot(HaveOccurred())

			By("checking that the logs are filtered correctly")
			Expect(page.Locator("table tr")).To(HaveCount(2))
			Expect(page.Locator("table tr:first-child td.message")).To(HaveText(lines[2]))
			Expect(page).ToNot(HaveText(lines[1]))

			_, _ = fmt.Fprintln(stdinWriter, lines[3])
			time.Sleep(100 * time.Millisecond)
			_, _ = fmt.Fprintln(stdinWriter, lines[4])

			By("checking that the logs are filtered correctly")
			Expect(page.Locator("table tr")).To(HaveCount(3))
			Expect(page.Locator("table tr:first-child td.message")).To(HaveText(lines[3]))
			Expect(page).ToNot(HaveText(lines[4]))

			By("clearing the filter")
			err = filterInput.Fill("")
			Expect(err).ShouldNot(HaveOccurred())

			By("checking that all logs are shown")
			Expect(page.Locator("table tr")).To(HaveCount(5))
			Expect(page.Locator("table tr:first-child td.message")).To(HaveText(lines[4]))
		})

		It("shows logs to multiple connected clients", func() {
			pw, err := playwright.Run()
			Expect(err).ShouldNot(HaveOccurred())

			browser, err := pw.Chromium.Launch()
			Expect(err).ShouldNot(HaveOccurred())

			var pages []playwright.Page

			for i := 0; i < 5; i++ {
				context, err := browser.NewContext()
				Expect(err).ShouldNot(HaveOccurred())

				page, err := context.NewPage()
				Expect(err).ShouldNot(HaveOccurred())

				pages = append(pages, page)

				_, err = page.Goto(targetUrl)
				Expect(err).ToNot(HaveOccurred())

				Expect(page).To(HaveText("Streamlog"))
			}

			By("giving it a second to have all the browsers and pages loaded")
			time.Sleep(1 * time.Second)

			By("sending a line to stdin")
			_, _ = fmt.Fprintln(stdinWriter, "and another")

			for i, page := range pages {
				By(fmt.Sprintf("checking page #%d", i+1))

				Expect(page).To(HaveText("and another"))
				Expect(page.Close()).ShouldNot(HaveOccurred())
			}
		})
	})
})
