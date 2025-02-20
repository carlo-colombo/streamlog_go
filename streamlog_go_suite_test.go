package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestStreamlogGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "StreamlogGo Suite")
}
