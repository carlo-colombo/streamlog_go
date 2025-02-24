package integration_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io"
	"os"
	"os/exec"
	"testing"
)

var (
	pathToBin                 string
	commitHash, pathToGitRepo string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)

	BeforeSuite(func() {
		var (
			err error
		)

		pathToBin, err = gexec.Build("github.com/carlo-colombo/streamlog_go")

		Expect(err).ToNot(HaveOccurred())
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
