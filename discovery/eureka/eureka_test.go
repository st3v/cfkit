package eureka

import (
	"time"

	"github.com/hudl/fargo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("eureka", func() {
	var (
		expectedURIs         = []string{"uri1", "uri2", "uri3"}
		expectedPort         = 12345
		expectedTimeout      = 987 * time.Second
		expectedPollInterval = 456 * time.Millisecond
	)

	Describe(".NewClient", func() {
		It("returns a client with a correctly initialized fargo connection", func() {
			c := NewClient(expectedURIs, expectedPort, expectedTimeout, expectedPollInterval)
			Expect(c.conn.ServiceUrls).To(HaveLen(len(expectedURIs)))
			Expect(c.conn.ServiceUrls).To(ConsistOf(expectedURIs))
			Expect(c.conn.ServicePort).To(Equal(expectedPort))
			Expect(c.conn.Timeout).To(Equal(expectedTimeout))
			Expect(c.conn.PollInterval).To(Equal(expectedPollInterval))
		})
	})

	Describe("fargoClient", func() {
		var client *fargoClient

		BeforeEach(func() {
			client = &fargoClient{
				conn: &fargo.EurekaConnection{
					ServiceUrls:  expectedURIs,
					ServicePort:  expectedPort,
					Timeout:      expectedTimeout,
					PollInterval: expectedPollInterval,
				},
			}
		})

		Describe(".URIs", func() {
			It("returns the service URIs of the underlying connection", func() {
				Expect(client.URIs()).To(HaveLen(len(expectedURIs)))
				Expect(client.URIs()).To(ConsistOf(expectedURIs))
			})

			Context("when the underlying connection is nil", func() {
				BeforeEach(func() {
					client = new(fargoClient)
				})

				It("returns an empty slice", func() {
					Expect(client.URIs()).To(HaveLen(0))
				})
			})
		})

		Describe(".Port", func() {
			It("returns the service port of the underlying connection", func() {
				Expect(client.Port()).To(Equal(expectedPort))
			})

			Context("when the underlying connection is nil", func() {
				BeforeEach(func() {
					client = new(fargoClient)
				})

				It("returns zero", func() {
					Expect(client.Port()).To(BeZero())
				})
			})
		})

		Describe(".Timeout", func() {
			It("returns the timeout for the underlying connection", func() {
				Expect(client.Timeout()).To(Equal(expectedTimeout))
			})

			Context("when the underlying connection is nil", func() {
				BeforeEach(func() {
					client = new(fargoClient)
				})

				It("returns zero", func() {
					Expect(client.Timeout()).To(BeZero())
				})
			})
		})

		Describe(".PollInterval", func() {
			It("returns the poll intervall for the underlying connection", func() {
				Expect(client.PollInterval()).To(Equal(expectedPollInterval))
			})

			Context("when the underlying connection is nil", func() {
				BeforeEach(func() {
					client = new(fargoClient)
				})

				It("returns zero", func() {
					Expect(client.PollInterval()).To(BeZero())
				})
			})
		})
	})
})
