package main_test

import (
	"fmt"
	"github.com/carlo-colombo/streamlog_go"
	"github.com/carlo-colombo/streamlog_go/logentry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
)

var _ = Describe("InMemoryLogsStore", func() {

	var (
		r     io.Reader
		w     *io.PipeWriter
		store *main.InMemoryLogsStore
	)

	BeforeEach(func() {
		r, w = io.Pipe()

		store = main.NewStore()

		go store.Scan(r)
	})

	It("scans and collect lines", func() {
		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
			_, _ = fmt.Fprintln(w, "New World")
		}()

		Eventually(store.List).Should(ContainElements(
			logentry.Log{Line: "Hello World"},
			logentry.Log{Line: "New World"},
		))
	})

	It("provide a channel that emits logs", func() {
		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
			_, _ = fmt.Fprintln(w, "New World")
		}()

		Eventually(store.LineFor("foo")).Should(Receive(
			Equal(logentry.Log{Line: "Hello World"})))
		Eventually(store.LineFor("foo")).Should(Receive(
			Equal(logentry.Log{Line: "New World"})))
	})

	It("support multiple clients consuming logs", func() {
		clientA := store.LineFor("client A")
		clientB := store.LineFor("client B")

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
		}()

		Eventually(clientA, "2s").Should(Receive(
			Equal(logentry.Log{Line: "Hello World"})))
		Eventually(clientB, "2s").Should(Receive(
			Equal(logentry.Log{Line: "Hello World"})))
	})

	It("removes client when disconnecting", func() {
		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
		}()

		Eventually(store.LineFor("client A")).Should(Receive(
			Equal(logentry.Log{Line: "Hello World"})))

		store.Disconnect("client A")

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World 2")
		}()

		Expect(store.Clients()).ToNot(ContainElements("client A"))
	})
})
