package main_test

import (
	"fmt"
	"io"

	main "github.com/carlo-colombo/streamlog_go"
	"github.com/carlo-colombo/streamlog_go/logentry"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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

		Eventually(store.List).Should(HaveLen(2))
		Eventually(store.List).Should(ContainElements(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World")),
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World")),
		))
	})

	It("provide a channel that emits logs", func() {
		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
			_, _ = fmt.Fprintln(w, "New World")
		}()

		Eventually(store.LineFor("foo")).Should(Receive(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World"))))
		Eventually(store.LineFor("foo")).Should(Receive(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World"))))
	})

	It("support multiple clients consuming logs", func() {
		clientA := store.LineFor("client A")
		clientB := store.LineFor("client B")

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
		}()

		Eventually(clientA, "2s").Should(Receive(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World"))))
		Eventually(clientB, "2s").Should(Receive(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World"))))
	})

	It("removes client when disconnecting", func() {
		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
		}()

		Eventually(store.LineFor("client A")).Should(Receive(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World"))))

		store.Disconnect("client A")

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World 2")
		}()

		Expect(store.Clients()).ToNot(ContainElements("client A"))
	})
})
