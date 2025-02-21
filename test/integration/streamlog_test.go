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
	"io"
	"net/http"
	"strings"
)

var _ = Describe("Test/Integration/Streamlog", func() {
	var session *gexec.Session
	var stdinReader io.Reader
	var stdinWriter *io.PipeWriter

	It("starts, and forward stdin to http response", func() {
		stdinReader, stdinWriter = io.Pipe()

		session = runBin([]string{}, stdinReader)
		Eventually(session.Err).Should(Say("Starting on http://localhost:"))

		targetUrl := getTargetUrl(session.Err)

		By(fmt.Sprintf("retrieving lines from endpoint %s", targetUrl))

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

	Describe("the endpoint", func() {
		It("returns JSON new line delimited body", func() {
			stdinReader, stdinWriter = io.Pipe()

			session = runBin([]string{}, stdinReader)

			Eventually(session.Err).Should(Say("Starting on http://localhost:"))

			targetUrl := getTargetUrl(session.Err)

			By(fmt.Sprintf("retrieving lines from endpoint %s", targetUrl))

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
			stdinReader, stdinWriter = io.Pipe()

			session = runBin([]string{}, stdinReader)

			Eventually(session.Err).Should(Say("Starting on http://localhost:"))

			targetUrl := getTargetUrl(session.Err)

			By(fmt.Sprintf("retrieving lines from endpoint %s", targetUrl))

			resp, err := http.Get(targetUrl + "/logs?sse")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(200))

			_, _ = fmt.Fprintln(stdinWriter, "and another")
			_, _ = fmt.Fprintln(stdinWriter, "line from stdin")

			scanner := bufio.NewScanner(resp.Body)
			scanner.Split(utils.ScanEvent)

			var events []string

			for scanner.Scan() {
				event := scanner.Text()
				events = append(events, event)

				if len(events) == 2 {
					resp.Body.Close()
					break
				}
			}

			Expect(events).To(ContainElements(
				`data: and another`, `data: line from stdin`,
			))
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
