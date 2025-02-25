package integration_test

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/onsi/gomega/types"
	"github.com/playwright-community/playwright-go"
	"io"
	"net/http"
	"strings"
	"time"
)

var expect playwright.PlaywrightAssertions = playwright.NewPlaywrightAssertions(2000)

var _ = Describe("Test/Integration/Streamlog", func() {
	var session *gexec.Session
	var stdinReader io.Reader
	var stdinWriter *io.PipeWriter
	var targetUrl string

	expect = playwright.NewPlaywrightAssertions()

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
				HaveHTTPHeaderWithValue("Content-Type", ContainSubstring("text/html")),
			))

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

	It("sends logs to a client that are ingested before the client is connected", func() {
		_, _ = fmt.Fprintln(stdinWriter, "and another")
		_, _ = fmt.Fprintln(stdinWriter, "line from stdin")

		for i := 0; i < 5; i++ {
			resp, err := http.Get(targetUrl + "/logs?sse")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp).To(SatisfyAll(
				HaveHTTPStatus(http.StatusOK),
				HaveHTTPHeaderWithValue("Content-Type", "text/event-stream"),
			))

			bodyReader := BufferReader(resp.Body)
			Eventually(bodyReader).Should(Say("and another"))
			Eventually(bodyReader).Should(Say("line from stdin"))
		}
	})

	AfterEach(func() {
		By("terminating the process")
		Expect(stdinWriter.Close()).ShouldNot(HaveOccurred())
		session.Terminate()
		Eventually(session).Should(gexec.Exit())
	})
})

type haveCount struct {
	expected int
	actual   int
}

func (h *haveCount) Match(actual interface{}) (success bool, err error) {
	actualLocator, ok := actual.(playwright.Locator)

	if !ok {
		return false, fmt.Errorf("HaveCount matcher expects a playwright.Locator")
	}

	err = expect.Locator(actualLocator).ToHaveCount(h.expected)
	locators, _ := actualLocator.All()
	h.actual = len(locators)

	return err == nil, nil
}

func (h haveCount) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(`Expected locator to have count of %d
Actual locator has count %d`, h.expected, h.actual)
}

func (h haveCount) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf(`Expected locator to not have count %d
Actual locator has count of %d`, h.expected, h.actual)
}

func HaveCount(expected int) types.GomegaMatcher {
	return &haveCount{expected: expected}
}

func HaveText(text string) types.GomegaMatcher {
	return &haveTextMatcher{text}
}

func HaveSelector(selector string) types.GomegaMatcher {
	return &haveSelectorMatcher{selector: selector}
}

type haveSelectorMatcher struct {
	selector string
}
type haveTextMatcher struct {
	text string
}

func (h haveTextMatcher) Match(actual interface{}) (success bool, err error) {
	if page, ok := actual.(playwright.Page); ok {
		err = expect.Locator(page.GetByText(h.text)).ToBeVisible()
		return err == nil, nil
	}

	if locator, ok := actual.(playwright.Locator); ok {
		err = expect.Locator(locator).ToHaveText(h.text)
		return err == nil, nil
	}
	return false, errors.New("HaveTextMatcher expects a playwright.Page or playwright.Locator")
}

func (h haveTextMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain text \n\t%#v", actual, h.text)
}

func (h haveTextMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\n not to contain text \n\t%#v", actual, h.text)
}

func (h haveSelectorMatcher) Match(actual interface{}) (success bool, err error) {
	page, ok := actual.(playwright.Page)

	if !ok {
		return false, fmt.Errorf("HaveSelector matcher expects a playwright.Page")
	}

	err = expect.Locator(page.Locator(h.selector)).ToBeVisible()

	return err == nil, err
}

func (h haveSelectorMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain selector \n\t%#v", actual, h.selector)
}

func (h haveSelectorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\n not to contain selector \n\t%#v", actual, h.selector)
}

func getTargetUrl(err *Buffer) string {
	targetUrl, _ := strings.CutPrefix(string(err.Contents()), "Starting on")
	targetUrl = strings.TrimSpace(targetUrl)
	return targetUrl
}
