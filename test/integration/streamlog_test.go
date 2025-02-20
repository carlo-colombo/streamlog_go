package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io"
	"net/http"
)

var _ = Describe("Test/Integration/Streamlog", func() {
	It("echoes multiple line from stdin", func() {
		r, w := io.Pipe()

		_, _, session := runBin([]string{}, r, 0)

		Eventually(session.Err).Should(Say("Starting on http://localhost:"))

		By("sending lines to stdin")

		fmt.Fprintln(w, "some line from stdin")
		Eventually(session).Should(Say("some line from stdin"))

		fmt.Fprintln(w, "and another")
		fmt.Fprintln(w, "line from stdin")
		Eventually(session).Should(Say("and another\nline from stdin"))

		By(fmt.Sprintf("retrieving lines from endpoint"))

		resp, err := http.Get("http://localhost:8080")
		Expect(err).ShouldNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(200))

		By("terminating the process")
		Expect(w.Close()).ShouldNot(HaveOccurred())
		session.Terminate()
		Eventually(session).Should(gexec.Exit())
	})
})
