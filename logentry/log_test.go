package logentry_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/carlo-colombo/streamlog_go/logentry"
)

var _ = Describe("Log", func() {
	Describe("Encode", func() {
		It("write the encoded log", func() {
			b := NewBuffer()
			e := json.NewEncoder(b)

			log := logentry.NewLog("message")
			log.Encode(e)
			b.Close()

			Eventually(b).Should(Say(`{"line":"message","timestamp":`))
		})

		It("handlers if the buffers is closed", func() {
			b := NewBuffer()
			e := json.NewEncoder(b)
			Expect(b.Close()).To(Succeed())

			log := logentry.NewLog("message")
			Expect(log.Encode(e)).
				To(MatchError(ContainSubstring("impossible to encode log entry:")))
		})
	})
})
