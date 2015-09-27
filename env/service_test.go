package env_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/cfkit/env"
)

var _ = Describe(".ServiceWithTag", func() {
	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServices)
	})

	AfterEach(func() {
		os.Unsetenv("VCAP_SERVICES")
	})

	Context("when the service is defined", func() {
		It("does NOT return an error", func() {
			_, err := env.ServiceWithTag("service-tag-2")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct service", func() {
			svc, _ := env.ServiceWithTag("service-tag-2")
			Expect(svc.Name).To(Equal("service-name-2"))
			Expect(svc.Label).To(Equal("service-label"))
			Expect(svc.Plan).To(Equal("service-plan-2"))
		})

		It("is case-insensitive", func() {
			_, err := env.ServiceWithTag("SERVICE-TAG-2")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when service is NOT defined", func() {
		It("returns an error", func() {
			_, err := env.ServiceWithTag("unknown")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Service with tag 'unknown' not found"))
		})
	})

	Context("when VCAP_SERVICES is not set", func() {
		BeforeEach(func() {
			os.Unsetenv("VCAP_SERVICES")
		})

		It("returns an error", func() {
			_, err := env.ServiceWithTag("service")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("VCAP_SERVICES not set"))
		})
	})

	Context("when VCAP_SERVICES can NOT be unmarshalled", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_SERVICES", "INVALID")
		})

		It("returns an error", func() {
			_, err := env.ServiceWithTag("service")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error parsing VCAP_SERVICES"))
		})
	})
})

var _ = Describe(".ServiceWithName", func() {
	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServices)
	})

	AfterEach(func() {
		os.Unsetenv("VCAP_SERVICES")
	})

	Context("when the service is defined", func() {
		It("does NOT return an error", func() {
			_, err := env.ServiceWithName("service-name-2")
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct service", func() {
			svc, _ := env.ServiceWithName("service-name-2")
			Expect(svc.Name).To(Equal("service-name-2"))
			Expect(svc.Label).To(Equal("service-label"))
			Expect(svc.Plan).To(Equal("service-plan-2"))
		})

		It("is case-insensitive", func() {
			_, err := env.ServiceWithName("SERVICE-NAME-2")
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when service is NOT defined", func() {
		It("returns an error", func() {
			_, err := env.ServiceWithName("unknown")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Service with name 'unknown' not found"))
		})
	})

	Context("when VCAP_SERVICES is not set", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_SERVICES", "")
		})

		It("returns an error", func() {
			_, err := env.ServiceWithName("service")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("VCAP_SERVICES not set"))
		})
	})

	Context("when VCAP_SERVICES can NOT be unmarshalled", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_SERVICES", "INVALID")
		})

		It("returns an error", func() {
			_, err := env.ServiceWithName("service")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error parsing VCAP_SERVICES"))
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
