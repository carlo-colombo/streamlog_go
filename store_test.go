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

		Eventually(store.List).Should(SatisfyAll(
			HaveLen(2),
			ContainElements(
				WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World")),
				WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World")),
			),
		))
	})

	It("filters logs based on case-insensitive search", func() {
		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
			_, _ = fmt.Fprintln(w, "New World")
			_, _ = fmt.Fprintln(w, "Another Line")
		}()

		Eventually(store.List).Should(HaveLen(3))

		go store.SetFilter("world")
		Eventually(store.List).Should(HaveLen(2))
		Eventually(store.List).Should(ContainElements(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World")),
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World")),
		))

		go store.SetFilter("")
		Eventually(store.List).Should(HaveLen(3))
	})

	It("provide a channel that emits logs", func() {
		client := store.LineFor("foo")

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
			_, _ = fmt.Fprintln(w, "New World")
		}()

		Eventually(client).Should(Receive(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World"))))
		Eventually(client).Should(Receive(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World"))))
	})

	It("support multiple clients consuming logs", func() {
		clientA := store.LineFor("client A")
		clientB := store.LineFor("client B")

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
		}()

		logsCh := make(chan string)

		// collect logs from both clients
		go func() {
			for {
				select {
				case l := <-clientA:
					logsCh <- l.Line
				case l := <-clientB:
					logsCh <- l.Line
				}
			}
		}()

		Eventually(logsCh).Should(Receive(Equal("Hello World")))
		Eventually(logsCh).Should(Receive(Equal("Hello World")))
	})

	It("removes client when disconnecting", func() {
		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
		}()

		Eventually(store.LineFor("client A")).Should(Receive(
			WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World"))))

		store.Disconnect("client A")

		Expect(store.Clients()).ToNot(ContainElements("client A"))
	})

	It("stores all logs regardless of filter", func() {
		// Set up filter change channel
		filterChangeCh := store.FilterChangeFor()

		// Set initial filter and wait for it to be applied
		go store.SetFilter("world")
		Eventually(filterChangeCh).Should(Receive())

		// Write logs after filter is set
		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
			_, _ = fmt.Fprintln(w, "Another Line")
			_, _ = fmt.Fprintln(w, "New World")
		}()

		// Wait for filtered results
		Eventually(func(g Gomega) {
			g.Expect(store.List()).To(SatisfyAll(
				HaveLen(2),
				ContainElements(
					WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World")),
					WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World")),
				),
			))
		}, "2s", "0.5s").Should(Succeed())

		// Clear filter and wait for it to be applied
		go store.SetFilter("")
		Eventually(filterChangeCh).Should(Receive())

		// Wait for unfiltered results
		Eventually(func(g Gomega) {
			g.Expect(store.List()).To(HaveLen(3))
			g.Expect(store.List()).To(ContainElements(
				WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World")),
				WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Another Line")),
				WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World")),
			))
		}, "2s").Should(Succeed())
	})

	It("emits a signal when filter changes", func() {
		// Get the filter change channel for a client
		filterChangeCh := store.FilterChangeFor()

		// Set initial filter
		go store.SetFilter("world")
		Eventually(filterChangeCh).Should(Receive())

		// Change filter to empty
		go store.SetFilter("")
		Eventually(filterChangeCh).Should(Receive())

		// Change filter to new value
		go store.SetFilter("test")
		Eventually(filterChangeCh).Should(Receive())
	})
})
