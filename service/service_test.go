package service_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/cfkit/service"
)

var _ = Describe("service", func() {

	Describe(".WithTag", func() {

		Context("when the service is defined", func() {

			BeforeEach(func() {
				jsonStr := `{
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
					}`
				os.Setenv("VCAP_SERVICES", jsonStr)
			})

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

		})

		Context("when service is NOT defined", func() {

			It("returns an error", func() {
				_, err := service.WithTag("not-there")
				Expect(err).To(HaveOccurred())
			})

		})

	})

})
