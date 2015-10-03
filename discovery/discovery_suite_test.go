package discovery_test

import (
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDiscovery(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Discovery Suite")
}
