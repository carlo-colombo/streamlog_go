package integration_test

import (
	"bufio"
	"fmt"
	"github.com/carlo-colombo/streamlog_go/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io"
	"net/http"
)

var _ = Describe("Test/Integration/Streamlog", func() {

	It("accepts port as parameter to expose the service", func() {
		stdinReader, stdinWriter = io.Pipe()

		session = runBin([]string{"--port", "32323"}, stdinReader)

		Eventually(session.Err).Should(Say("Starting on http://localhost:32323"))

		resp, err := http.Get(getTargetUrl(session.Err) + "/logs")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))
	})

	Describe("/logs endpoint", func() {
		It("streams events matching the lines read from stdin", func() {
			By("sending lines before connecting to the endpoint")

			_, _ = fmt.Fprintln(stdinWriter, "first line")
			_, _ = fmt.Fprintln(stdinWriter, "second line")

			By("requesting the logs endpoint")

			resp, err := http.Get(targetUrl + "/logs")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp).To(SatisfyAll(
				HaveHTTPStatus(http.StatusOK),
				HaveHTTPHeaderWithValue("Content-Type", "text/event-stream"),
			))

			By("checking the response")

			scanner := bufio.NewScanner(resp.Body)
			scanner.Split(utils.ScanEvent)

			Expect(scanner.Scan()).To(BeTrue())
			Expect(scanner.Text()).To(MatchRegexp("data:.*first line"))

			Expect(scanner.Scan()).To(BeTrue())
			Expect(scanner.Text()).To(MatchRegexp("data:.*second line"))

			By("sending lines after connecting to the endpoint")

			_, _ = fmt.Fprintln(stdinWriter, "and another")
			_, _ = fmt.Fprintln(stdinWriter, "line from stdin")

			Expect(scanner.Scan()).To(BeTrue())
			Expect(scanner.Text()).To(MatchRegexp("data:.*and another"))

			Expect(scanner.Scan()).To(BeTrue())
			Expect(scanner.Text()).To(MatchRegexp("data:.*line from stdin"))
		})
	})

	Describe("/clients endpoint", func() {
		It("returns a count of clients", func() {
			Expect(http.Get(targetUrl + "/clients")).To(
				SatisfyAll(
					HaveHTTPStatus(http.StatusOK),
					HaveHTTPBody("0")))

			resp, _ := http.Get(targetUrl + "/logs?sse")

			Expect(http.Get(targetUrl + "/clients")).To(
				SatisfyAll(
					HaveHTTPStatus(http.StatusOK),
					HaveHTTPBody("1")))

			By("having the client closing the connection")
			Expect(resp.Body.Close()).ToNot(HaveOccurred())

			Eventually(func() (*http.Response, error) {
				return http.Get(targetUrl + "/clients")
			}).Should(
				SatisfyAll(
					HaveHTTPStatus(http.StatusOK),
					HaveHTTPBody("0")))
		})
	})

	AfterEach(func() {
		By("terminating the process")
		Expect(stdinWriter.Close()).ShouldNot(HaveOccurred())
		session.Terminate()
		Eventually(session).Should(gexec.Exit())
	})
})
