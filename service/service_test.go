package service_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/cfkit/service"
)

var _ = Describe("service", func() {
	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServices)
	})

	Describe(".WithTag", func() {
		Context("when the service is defined", func() {
			It("does NOT return an error", func() {
				_, err := service.WithTag("service-tag-2")
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns the correct service", func() {
				svc, _ := service.WithTag("service-tag-2")
				Expect(svc.Name).To(Equal("service-name-2"))
				Expect(svc.Label).To(Equal("service-label"))
				Expect(svc.Plan).To(Equal("service-plan-2"))
			})

			It("is case-insensitive", func() {
				_, err := service.WithTag("SERVICE-TAG-2")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when service is NOT defined", func() {
			It("returns an error", func() {
				_, err := service.WithTag("unknown")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Service with tag 'unknown' not found"))
			})
		})
	})

	Describe(".WithName", func() {
		Context("when the service is defined", func() {
			It("does NOT return an error", func() {
				_, err := service.WithName("service-name-2")
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns the correct service", func() {
				svc, _ := service.WithName("service-name-2")
				Expect(svc.Name).To(Equal("service-name-2"))
				Expect(svc.Label).To(Equal("service-label"))
				Expect(svc.Plan).To(Equal("service-plan-2"))
			})

			It("is case-insensitive", func() {
				_, err := service.WithName("SERVICE-NAME-2")
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when service is NOT defined", func() {
			It("returns an error", func() {
				_, err := service.WithName("unknown")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Service with name 'unknown' not found"))
			})
		})
	})
})

var vcapServices = `
	{
		"service-label": [
			{
				"name": "service-name-1",
				"label": "service-label",
				"tags": [ "foo", "bar", "service-tag-1" ],
				"plan": "service-plan-1",
				"credentials": {
					"username": "some-username",
					"password": "some-password"
				}
			},
			{
				"name": "service-name-2",
				"label": "service-label",
				"tags": [ "foo", "bar", "service-tag-2" ],
				"plan": "service-plan-2",
				"credentials": {
					"username": "some-username",
					"password": "some-password"
				}
			}
		]
	}
`
