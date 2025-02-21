package utils_test

import (
	"bufio"
	"github.com/carlo-colombo/streamlog_go/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strings"
)

var _ = Describe("Test/ScanSse", func() {
	sseStream := `event: foobar
data: cool cool cool

data: new cool cool

`
	It("scan_sse", func() {
		scanner := bufio.NewScanner(strings.NewReader(sseStream))
		scanner.Split(utils.ScanEvent)

		var events []string
		for scanner.Scan() {
			event := scanner.Text()
			events = append(events, event)

			if len(events) == 2 {
				break
			}
		}

		Expect(events).Should(ContainElements("event: foobar\ndata: cool cool cool", "data: new cool cool"))
	})

	It("Scan events", func() {
		advance, token, err := utils.ScanEvent([]byte(sseStream), false)

		Expect(err).ToNot(HaveOccurred())
		Expect(string(token)).Should(Equal("event: foobar\ndata: cool cool cool"))
		Expect(advance).Should(BeNumerically("==", 36))

		advance, token, err = utils.ScanEvent([]byte(sseStream[advance:]), false)

		Expect(err).ToNot(HaveOccurred())
		Expect(string(token)).Should(Equal("data: new cool cool"))
		Expect(advance).Should(BeNumerically("==", 21))
	})
})
