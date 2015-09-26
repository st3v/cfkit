package rabbitmq_test

import (
	"github.com/st3v/cfkit/service"
	"github.com/st3v/cfkit/service/rabbitmq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(".URI", func() {
	var srv service.Service

	Context("when the service is labeled p-rabbitmq", func() {
		BeforeEach(func() {
			protos := map[string]interface{}{
				"amqp": map[string]interface{}{
					"uris": []interface{}{"uri0", "uri2", "uri3"},
				},
			}

			srv = service.Service{
				Label:       "p-rabbitmq",
				Credentials: map[string]interface{}{"protocols": protos},
			}
		})

		It("does not return an error", func() {
			_, err := rabbitmq.URI(srv)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the first amqp URI for the service", func() {
			uri, _ := rabbitmq.URI(srv)
			Expect(uri).To(Equal("uri0"))
		})

		Context("when service protocols field is missing", func() {
			BeforeEach(func() {
				srv.Credentials = map[string]interface{}{}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid service credentials"))
			})
		})

		Context("when protocols field is of an invalid type", func() {
			BeforeEach(func() {
				srv.Credentials = map[string]interface{}{
					"protocols": false,
				}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid service credentials"))
			})
		})

		Context("when amqp protocol is missing", func() {
			BeforeEach(func() {
				protos := map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": false,
					},
				}

				srv.Credentials = map[string]interface{}{
					"protocols": protos,
				}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid AMQP protocol credentials"))
			})
		})

		Context("when the amqp field is of an invalid type", func() {
			BeforeEach(func() {
				protos := map[string]interface{}{
					"amqp": false,
				}

				srv.Credentials = map[string]interface{}{
					"protocols": protos,
				}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid AMQP protocol credentials"))
			})
		})

		Context("when amqp uris field is missing", func() {
			BeforeEach(func() {
				protos := map[string]interface{}{
					"amqp": map[string]interface{}{
						"bar": false,
					},
				}

				srv.Credentials = map[string]interface{}{
					"protocols": protos,
				}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid AMQP URIs"))
			})
		})

		Context("when amqp uris array is empty", func() {
			BeforeEach(func() {
				protos := map[string]interface{}{
					"amqp": map[string]interface{}{
						"uris": []interface{}{},
					},
				}

				srv.Credentials = map[string]interface{}{
					"protocols": protos,
				}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid AMQP URIs"))
			})
		})

		Context("when amqp uri is an empty string", func() {
			BeforeEach(func() {
				protos := map[string]interface{}{
					"amqp": map[string]interface{}{
						"uris": []interface{}{""},
					},
				}

				srv.Credentials = map[string]interface{}{
					"protocols": protos,
				}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid AMQP URI"))
			})
		})
	})

	Context("when the service is NOT labeled p-rabbitmq", func() {
		BeforeEach(func() {
			srv = service.Service{
				Label: "cloudamqp-dev",
				Credentials: map[string]interface{}{
					"uri": "amqp://uri",
				},
			}
		})

		It("does not return an error", func() {
			_, err := rabbitmq.URI(srv)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the correct URI", func() {
			uri, _ := rabbitmq.URI(srv)
			Expect(uri).To(Equal("amqp://uri"))
		})

		Context("when service credentials do NOT contain URI", func() {
			BeforeEach(func() {
				srv.Credentials = map[string]interface{}{}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid AMQP URI"))
			})
		})

		Context("when service URI is empty", func() {
			BeforeEach(func() {
				srv.Credentials = map[string]interface{}{
					"uri": "",
				}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid AMQP URI"))
			})
		})

		Context("when service URI does not specify amqp protocol", func() {
			BeforeEach(func() {
				srv.Credentials = map[string]interface{}{
					"uri": "http://foo.bar",
				}
			})

			It("returns expected error", func() {
				_, err := rabbitmq.URI(srv)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Invalid AMQP URI"))
			})
		})
	})
})
