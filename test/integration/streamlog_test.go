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
	It("echoes multiple line from stdin", func() {
		r, w := io.Pipe()

		session := runBin([]string{}, r)

		Eventually(session.Err).Should(Say("Starting on http://localhost:"))

		targetUrl := getTargetUrl(session.Err)

		By(fmt.Sprintf("retrieving lines from endpoint %s", targetUrl))

		resp, err := http.Get(targetUrl)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))

		By("sending lines to stdin and checking stdout")

		fmt.Fprintln(w, "some line from stdin")
		Eventually(session).Should(Say("some line from stdin"))

		fmt.Fprintln(w, "and another")
		fmt.Fprintln(w, "line from stdin")
		Eventually(session).Should(Say("and another\nline from stdin"))

		By("checking the response from the endpoint")
		Eventually(BufferReader(resp.Body)).Should(Say("some line from stdin"))

		By("terminating the process")
		Expect(w.Close()).ShouldNot(HaveOccurred())
		session.Terminate()
		Eventually(session).Should(gexec.Exit())
	})

	It("accepts port as parameter", func() {
		r, _ := io.Pipe()

		session := runBin([]string{"--port", "32323"}, r)

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
