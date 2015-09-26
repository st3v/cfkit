package rabbitmq_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRabbitmq(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Rabbitmq Suite")
}
