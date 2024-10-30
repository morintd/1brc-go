package test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBRC(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "1BRC Suite")
}
