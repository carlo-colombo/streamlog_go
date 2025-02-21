package log_test

import (
	"encoding/json"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"

	"github.com/carlo-colombo/streamlog_go/log"
)

var _ = Describe("Log", func() {
	Describe("Encode", func() {
		It("write the encoded log", func() {
			b := NewBuffer()
			e := json.NewEncoder(b)

			log.Log{Line: "message"}.Encode(e)
			b.Close()

			Eventually(b).Should(Say("message"))
		})
	})
})
