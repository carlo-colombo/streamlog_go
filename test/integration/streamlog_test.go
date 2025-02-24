package integration_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/playwright-community/playwright-go"
	"io"
	"net/http"
	"strings"
	"time"
)

var _ = Describe("Test/Integration/Streamlog", func() {
	var session *gexec.Session
	var stdinReader io.Reader
	var stdinWriter *io.PipeWriter
	var targetUrl string

	BeforeEach(func() {
		stdinReader, stdinWriter = io.Pipe()

		session = runBin([]string{}, stdinReader)

		Eventually(session.Err, "2s").Should(Say("Starting on http://localhost:"))

		targetUrl = getTargetUrl(session.Err)

		By(fmt.Sprintf("targetting endpoint %s", targetUrl))
	})

	It("forwards stdin to an endpoint", func() {
		resp, err := http.Get(targetUrl + "/logs")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))

		By("sending lines to stdin and checking stdout")

		_, _ = fmt.Fprintln(stdinWriter, "some line from stdin")

		By("checking the response from the endpoint")
		bodyReader := BufferReader(resp.Body)
		Eventually(bodyReader).Should(Say("some line from stdin"))

		By("sending multiple lines to stdin and checking the response from the endpoint")
		_, _ = fmt.Fprintln(stdinWriter, "and another")
		_, _ = fmt.Fprintln(stdinWriter, "line from stdin")

		Eventually(bodyReader).Should(Say("and another"))
		Eventually(bodyReader).Should(Say("line from stdin"))

		By("terminating the process")
		Expect(stdinWriter.Close()).ShouldNot(HaveOccurred())
		session.Terminate()
		Eventually(session).Should(gexec.Exit())
	})

	It("accepts port as parameter", func() {
		stdinReader, stdinWriter = io.Pipe()

		session = runBin([]string{"--port", "32323"}, stdinReader)

		Eventually(session.Err).Should(Say("Starting on http://localhost:32323"))

		resp, err := http.Get(getTargetUrl(session.Err) + "/logs")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))
	})

	Describe("the logs endpoint", func() {
		It("returns JSON new line delimited body", func() {
			resp, err := http.Get(targetUrl + "/logs")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))

			_, _ = fmt.Fprintln(stdinWriter, "and another")
			_, _ = fmt.Fprintln(stdinWriter, "line from stdin")

			scanner := bufio.NewScanner(resp.Body)

			i := 0
			for scanner.Scan() {
				i++
				result := make(map[string]interface{})
				lineBuffer := bytes.NewBuffer(scanner.Bytes())
				err := json.NewDecoder(lineBuffer).Decode(&result)
				Expect(err).ShouldNot(HaveOccurred())

				if i == 2 {
					_ = resp.Body.Close()
				}
			}
		})

		It("returns sse events with html content", func() {
			resp, err := http.Get(targetUrl + "/logs?sse")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp).To(SatisfyAll(
				HaveHTTPStatus(http.StatusOK),
				HaveHTTPHeaderWithValue("Content-Type", "text/event-stream"),
			))

			_, _ = fmt.Fprintln(stdinWriter, "and another")
			_, _ = fmt.Fprintln(stdinWriter, "line from stdin")

			scanner := bufio.NewScanner(resp.Body)
			scanner.Split(utils.ScanEvent)

			var events []string

			for scanner.Scan() {
				event := scanner.Text()
				events = append(events, event)

				if len(events) == 2 {
					_ = resp.Body.Close()
					break
				}
			}

			Expect(events).To(ContainElements(
				ContainSubstring(`and another`),
				ContainSubstring(`line from stdin`),
			))
		})
	})

	Describe("the root endpoint", func() {
		BeforeEach(func() { PauseOutputInterception() })

		AfterEach(func() { ResumeOutputInterception() })

		It("returns an index.html page", func() {
			resp, err := http.Get(targetUrl + "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(resp).To(SatisfyAll(
				HaveHTTPStatus(http.StatusOK),
				HaveHTTPHeaderWithValue("Content-Type", "text/html"),
			))

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

			_, _ = fmt.Fprintln(stdinWriter, "and another")
			_, _ = fmt.Fprintln(stdinWriter, "line from stdin")

			entries, err := page.Locator("table tr").All()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(entries).To(HaveLen(2))

			Expect(entries[0].TextContent()).To(ContainSubstring("line from stdin"))
			Expect(entries[1].TextContent()).To(ContainSubstring("and another"))

			By("sending an additional line to stdin and prepending to the content")
			_, _ = fmt.Fprintln(stdinWriter, "bonus line from stdin")
			logLine, err := page.Locator("table tr:first-child").AllTextContents()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(logLine[0]).To(ContainSubstring("bonus line from stdin"))
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

				_ = page.GetByText("Streamlog")
			}

			By("giving it a second to have all the browsers and pages loaded")
			time.Sleep(1 * time.Second)

			By("sending a line to stdin")
			_, _ = fmt.Fprintln(stdinWriter, "and another")

			expect := playwright.NewPlaywrightAssertions()

			for i, page := range pages {
				By(fmt.Sprintf("checking page #%d", i+1))

				err := expect.Locator(page.GetByText("and another")).ToBeVisible()
				Expect(err).ShouldNot(HaveOccurred())

				Expect(page.Close()).ShouldNot(HaveOccurred())
			}
		})
	})

	AfterEach(func() {
		By("terminating the process")
		Expect(stdinWriter.Close()).ShouldNot(HaveOccurred())
		session.Terminate()
		Eventually(session).Should(gexec.Exit())
	})
})

func getTargetUrl(err *Buffer) string {
	targetUrl, _ := strings.CutPrefix(string(err.Contents()), "Starting on")
	targetUrl = strings.TrimSpace(targetUrl)
	return targetUrl
}
