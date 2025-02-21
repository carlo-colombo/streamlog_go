package sse_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sse Suite")
}
