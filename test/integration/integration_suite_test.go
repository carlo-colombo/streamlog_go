package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/playwright-community/playwright-go"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var (
	pathToBin     string
	pathToGitRepo string
)

var (
	session     *gexec.Session
	stdinReader io.Reader
	stdinWriter *io.PipeWriter
	targetUrl   string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var (
			err error
		)

		pathToBin, err = gexec.Build("github.com/carlo-colombo/streamlog_go")

		Expect(err).ToNot(HaveOccurred())
		PauseOutputInterception()
		Expect(playwright.Install()).ToNot(HaveOccurred())
		ResumeOutputInterception()
	})

	BeforeEach(func() {
		stdinReader, stdinWriter = io.Pipe()

		session = runBin([]string{}, stdinReader)

		Eventually(session.Err, "2s").Should(Say("Starting on http://localhost:"))

		targetUrl = getTargetUrl(session.Err)

		By(fmt.Sprintf("targetting endpoint %s", targetUrl))
	})

	AfterSuite(func() {
		os.RemoveAll(pathToGitRepo)
		gexec.Kill()
		gexec.CleanupBuildArtifacts()
	})

	RunSpecs(t, "Integration Suite")
}

func runBin(args []string, stdIn io.Reader) (session *gexec.Session) {
	cmd := exec.Command(pathToBin, args...)
	cmd.Stdin = stdIn
	session, err := gexec.Start(cmd,
		gexec.NewPrefixedWriter("[streamlog out] ", GinkgoWriter),
		gexec.NewPrefixedWriter("[streamlog err] ", GinkgoWriter))
	Expect(err).ToNot(HaveOccurred())

	return session
}

func getTargetUrl(err *Buffer) string {
	targetUrl, _ := strings.CutPrefix(string(err.Contents()), "Starting on")
	targetUrl = strings.TrimSpace(targetUrl)
	return targetUrl
}
