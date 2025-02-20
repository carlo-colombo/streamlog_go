package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io"
	"net/http"
	"strings"
)

var _ = Describe("Test/Integration/Streamlog", func() {
	It("starts, and forward stdin to http response", func() {
		stdinReader, stdinWriter := io.Pipe()

		session := runBin([]string{}, stdinReader)

		Eventually(session.Err).Should(Say("Starting on http://localhost:"))

		targetUrl := getTargetUrl(session.Err)

		By(fmt.Sprintf("retrieving lines from endpoint %s", targetUrl))

		resp, err := http.Get(targetUrl)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))

		By("sending lines to stdin and checking stdout")

		fmt.Fprintln(stdinWriter, "some line from stdin")

		By("checking the response from the endpoint")
		bodyReader := BufferReader(resp.Body)
		Eventually(bodyReader).Should(Say("some line from stdin"))

		By("sending multiple lines to stdin and checking the response from the endpoint")
		fmt.Fprintln(stdinWriter, "and another")
		fmt.Fprintln(stdinWriter, "line from stdin")

		Eventually(bodyReader).Should(Say("and another"))
		Eventually(bodyReader).Should(Say("line from stdin"))

		By("terminating the process")
		Expect(stdinWriter.Close()).ShouldNot(HaveOccurred())
		session.Terminate()
		Eventually(session).Should(gexec.Exit())
	})

	It("accepts port as parameter", func() {
		stdinReader, _ := io.Pipe()

		session := runBin([]string{"--port", "32323"}, stdinReader)

		Eventually(session.Err).Should(Say("Starting on http://localhost:32323"))

		resp, err := http.Get(getTargetUrl(session.Err))
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))
	})
})

func getTargetUrl(err *Buffer) string {
	targetUrl, _ := strings.CutPrefix(string(err.Contents()), "Starting on")
	targetUrl = strings.TrimSpace(targetUrl)
	return targetUrl
}
