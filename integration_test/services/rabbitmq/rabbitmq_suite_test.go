package rabbitmq_test

import (
	"log"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRabbitmq(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())

	log.SetOutput(GinkgoWriter)

	RegisterFailHandler(Fail)
	RunSpecs(t, "RabbitMQ Integration Suite")
}
