package sse_test

import (
	"github.com/carlo-colombo/streamlog_go/logentry"
	"github.com/carlo-colombo/streamlog_go/sse"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Sse/Encoder", func() {

	It("encodes a single event", func() {
		buffer := gbytes.NewBuffer()
		e := sse.NewEncoder(buffer)

		Expect(e).ToNot(BeNil())

		err := e.Encode(logentry.Log{Line: "foobar"})
		Expect(err).ToNot(HaveOccurred())

		Eventually(buffer).Should(gbytes.Say("data: {\"line\":\"foobar\",\"timestamp\":\"0001-01-01T00:00:00Z\"}\n\n"))
	})

	It("returns an error if is not a log", func() {
		e := sse.NewEncoder(gbytes.NewBuffer())

		Expect(e).ToNot(BeNil())

		err := e.Encode(nil)
		Expect(err).To(MatchError(ContainSubstring("encoder can only encode a log object")))
	})
})
