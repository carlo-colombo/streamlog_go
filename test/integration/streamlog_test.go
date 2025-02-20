package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io"
)

var _ = Describe("Test/Integration/Streamlog", func() {
	It("echoes multiple line from stdin", func() {
		r, w := io.Pipe()

		_, _, session := runBin([]string{}, r, 0)

		fmt.Fprintln(w, "some line from stdin")

		Eventually(session).Should(Say("some line from stdin"))

		fmt.Fprintln(w, "and another")
		fmt.Fprintln(w, "line from stdin")

		Eventually(session).Should(Say("and another\nline from stdin"))

		Expect(w.Close()).ToNot(HaveOccurred())

		session.Terminate()

		Eventually(session).Should(gexec.Exit())
	})
})
