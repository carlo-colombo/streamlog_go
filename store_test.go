package main_test

import (
	"fmt"
	"github.com/carlo-colombo/streamlog_go/logentry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"

	"github.com/carlo-colombo/streamlog_go"
)

var _ = Describe("Store", func() {

	It("scans and collect lines", func() {
		var r, w = io.Pipe()

		store := main.NewStore()

		go store.Scan(r)

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
		var r, w = io.Pipe()

		store := main.NewStore()

		go store.Scan(r)

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
		var r, w = io.Pipe()

		store := main.NewStore()

		go store.Scan(r)

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
		}()

		Eventually(store.LineFor("client A")).Should(Receive(
			Equal(logentry.Log{Line: "Hello World"})))
		Eventually(store.LineFor("client B")).Should(Receive(
			Equal(logentry.Log{Line: "Hello World"})))
	})

	It("removes client when unsubscribing", func() {
		var r, w = io.Pipe()

		store := main.NewStore()

		go store.Scan(r)

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
		}()

		Eventually(store.LineFor("client A")).Should(Receive(
			Equal(logentry.Log{Line: "Hello World"})))

		store.Unsubscribe("client A")

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World 2")
		}()

		Expect(store.Clients()).ToNot(ContainElements("client A"))
	})
})
