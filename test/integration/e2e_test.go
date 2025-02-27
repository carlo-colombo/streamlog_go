package integration_test

import (
	"fmt"
	. "github.com/carlo-colombo/streamlog_go/test/playwrightmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/playwright-community/playwright-go"
	"net/http"
	"time"
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

			browser, err := pw.Chromium.Launch()
			Expect(err).ShouldNot(HaveOccurred())

			page, err := browser.NewPage()
			Expect(err).ShouldNot(HaveOccurred())

			By("opening the web browser")
			_, err = page.Goto(targetUrl)
			Expect(err).ShouldNot(HaveOccurred())

			Expect(page).To(HaveSelector("h1"))
			Expect(page).To(HaveText("Streamlog"))

			_, _ = fmt.Fprintln(stdinWriter, "line from stdin")
			_, _ = fmt.Fprintln(stdinWriter, "and another")

			Expect(page).To(HaveText("line from stdin"))
			Expect(page).To(HaveText("and another"))

			Expect(page.Locator("table tr")).To(HaveCount(2))

			By("sending an additional line to stdin and prepending to the content")
			_, _ = fmt.Fprintln(stdinWriter, "bonus line from stdin")
			Expect(page.Locator("table tr:first-child")).To(HaveText("bonus line from stdin"))
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
