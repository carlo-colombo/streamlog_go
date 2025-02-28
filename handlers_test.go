package main_test

import (
	"bufio"
	"context"
	"github.com/carlo-colombo/streamlog_go"
	"github.com/carlo-colombo/streamlog_go/logentry"
	"github.com/carlo-colombo/streamlog_go/test/utils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
)

type mockStore struct {
	clients      []string
	logs         []string
	logsCh       chan logentry.Log
	disconnected bool
}

func (m *mockStore) Scan(r io.Reader) {
	panic("implement me")
}

func (m *mockStore) List() []logentry.Log {
	var l []logentry.Log

	for _, line := range m.logs {
		l = append(l, logentry.Log{Line: line})
	}
	return l
}

func (m *mockStore) Disconnect(uid string) {
	m.disconnected = true
}

func (m *mockStore) LineFor(uid string) chan logentry.Log {
	return m.logsCh
}

func (m *mockStore) Clients() []string {
	return m.clients
}

var _ = Describe("Handlers", func() {
	var req *http.Request
	var rr *httptest.ResponseRecorder

	BeforeEach(func() {
		var err error

		req, err = http.NewRequest("", "", nil)
		Expect(err).NotTo(HaveOccurred())

		rr = httptest.NewRecorder()
	})

	Describe("ClientsHandler", func() {
		It("writes 0 on the response when no client is attached", func() {
			var store = &mockStore{}
			clientsHandlerFunc := main.ClientsHandler(store)

			handler := http.HandlerFunc(clientsHandlerFunc)

			handler.ServeHTTP(rr, req)

			Expect(rr).To(SatisfyAll(
				HaveHTTPStatus(http.StatusOK),
				HaveHTTPBody(Equal([]byte("0"))),
			))
		})

		It("returns the updated count of attached clients", func() {
			var store = &mockStore{clients: []string{"client1", "client2"}}
			clientsHandlerFunc := main.ClientsHandler(store)

			handler := http.HandlerFunc(clientsHandlerFunc)

			handler.ServeHTTP(rr, req)

			Expect(rr).To(SatisfyAll(
				HaveHTTPStatus(http.StatusOK),
				HaveHTTPBody(Equal([]byte("2"))),
			))

			store.clients = store.clients[1:]

			rr = httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			Expect(rr).To(SatisfyAll(
				HaveHTTPStatus(http.StatusOK),
				HaveHTTPBody(Equal([]byte("1"))),
			))
		})
	})

	Describe("logs handler", func() {
		It("writes the collected logs", func() {
			var store = &mockStore{logs: []string{"log1", "log2"}}
			clientsHandlerFunc := main.LogsHandler(store)

			handler := http.HandlerFunc(clientsHandlerFunc)

			go func() {
				handler.ServeHTTP(rr, req)
			}()

			Expect(rr).To(HaveHTTPStatus(http.StatusOK))

			Eventually(func(g Gomega) {
				scanner := bufio.NewScanner(rr.Body)
				scanner.Split(utils.ScanEvent)

				g.Expect(scanner.Scan()).To(BeTrue())
				g.Expect(scanner.Text()).To(MatchRegexp("data:.*log1"))

				g.Expect(scanner.Scan()).To(BeTrue())
				g.Expect(scanner.Text()).To(MatchRegexp("data:.*log2"))
			}).Should(Succeed())
		})

		It("streams additional logs", func() {
			var store = &mockStore{logsCh: make(chan logentry.Log)}
			clientsHandlerFunc := main.LogsHandler(store)

			handler := http.HandlerFunc(clientsHandlerFunc)

			go func() {
				handler.ServeHTTP(rr, req)
			}()

			go func() {
				store.logsCh <- logentry.Log{Line: "log1"}
			}()

			Expect(rr).To(HaveHTTPStatus(http.StatusOK))

			Eventually(func(g Gomega) {
				scanner := bufio.NewScanner(rr.Body)
				scanner.Split(utils.ScanEvent)

				g.Expect(scanner.Scan()).To(BeTrue())
				g.Expect(scanner.Text()).To(MatchRegexp("data:.*log1"))
			}).Should(Succeed())
		})

		It("disconnects clients when the client closes the connection", func() {
			var store = &mockStore{logs: []string{"log1", "log2"}}
			clientsHandlerFunc := main.LogsHandler(store)

			handler := http.HandlerFunc(clientsHandlerFunc)

			closeConnectionCtx, closeConnectionFunc := context.WithCancel(req.Context())

			go func() {
				handler.ServeHTTP(rr, req.WithContext(closeConnectionCtx))
			}()

			Expect(rr).To(HaveHTTPStatus(http.StatusOK))

			closeConnectionFunc()

			Eventually(func(g Gomega) {
				g.Expect(store.disconnected).To(BeTrue())
			}, "2s").Should(Succeed())
		})
	})
})
