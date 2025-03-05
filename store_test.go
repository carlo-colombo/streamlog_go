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
		store.SetFilter("world")

		go func() {
			_, _ = fmt.Fprintln(w, "Hello World")
			_, _ = fmt.Fprintln(w, "Another Line")
			_, _ = fmt.Fprintln(w, "New World")
		}()

		Eventually(func(g Gomega) {
			g.Expect(store.List()).To(SatisfyAll(
				HaveLen(2),
				ContainElements(
					WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World")),
					WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World")),
				),
			))
		}).Should(Succeed())

		store.SetFilter("")

		Eventually(func(g Gomega) {
			g.Expect(store.List()).To(HaveLen(3))
			g.Expect(store.List()).To(ContainElements(
				WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Hello World")),
				WithTransform(func(l logentry.Log) string { return l.Line }, Equal("Another Line")),
				WithTransform(func(l logentry.Log) string { return l.Line }, Equal("New World")),
			))
		}).Should(Succeed())
	})

	It("emits a signal when filter changes only if clients are connected", func() {
		// Get the filter change channel for a client
		filterChangeCh := store.FilterChangeFor()

		// Set filter with no clients connected - should not emit signal
		go store.SetFilter("world")
		Consistently(filterChangeCh).ShouldNot(Receive())

		// Connect a client by getting their line channel
		store.LineFor("client-1")

		// Now set filter - should emit signal since client is connected
		go store.SetFilter("test")
		Eventually(filterChangeCh).Should(Receive())

		// Disconnect the client
		store.Disconnect("client-1")

		// Set filter again with no clients - should not emit signal
		go store.SetFilter("another")
		Consistently(filterChangeCh).ShouldNot(Receive())
	})
})
