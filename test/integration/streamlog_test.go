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

		_, _, session := runBin([]string{}, r, 0)

		Eventually(session.Err).Should(Say("Starting on http://localhost:"))

		targetUrl, _ := strings.CutPrefix(string(session.Err.Contents()), "Starting on")
		targetUrl = strings.TrimSpace(targetUrl)

		resp, err := http.Get(targetUrl)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))

		By("sending lines to stdin and checking stdout")

		fmt.Fprintln(w, "some line from stdin")
		Eventually(session).Should(Say("some line from stdin"))

		fmt.Fprintln(w, "and another")
		fmt.Fprintln(w, "line from stdin")
		Eventually(session).Should(Say("and another\nline from stdin"))

		By(fmt.Sprintf("retrieving lines from endpoint %s", targetUrl))

		By("checking the response from the endpoint")
		Eventually(BufferReader(resp.Body)).Should(Say("some line from stdin"))

		By("terminating the process")
		Expect(w.Close()).ShouldNot(HaveOccurred())
		session.Terminate()
		Eventually(session).Should(gexec.Exit())
	})
})
