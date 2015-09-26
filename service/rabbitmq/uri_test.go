package rabbitmq_test

import (
	"github.com/st3v/cfkit/service"
	"github.com/st3v/cfkit/service/rabbitmq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe(".URI", func() {
	var srv service.Service

	BeforeEach(func() {
		protos := map[string]interface{}{
			"amqp": map[string]interface{}{
				"uris": []interface{}{"uri0", "uri2", "uri3"},
			},
		}

		srv = service.Service{
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
		var svc = service.Service{}

		It("returns expected error", func() {
			_, err := rabbitmq.URI(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid service credentials"))
		})
	})

	Context("when protocols field is of an invalid type", func() {
		var svc = service.Service{
			Credentials: map[string]interface{}{
				"protocols": false,
			},
		}

		It("returns expected error", func() {
			_, err := rabbitmq.URI(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid service credentials"))
		})
	})

	Context("when amqp protocol is missing", func() {
		var (
			protos = map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": false,
				},
			}

			svc = service.Service{
				Credentials: map[string]interface{}{
					"protocols": protos,
				},
			}
		)

		It("returns expected error", func() {
			_, err := rabbitmq.URI(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid AMQP protocol credentials"))
		})
	})

	Context("when the amqp field is of an invalid type", func() {
		var (
			protos = map[string]interface{}{
				"amqp": false,
			}

			svc = service.Service{
				Credentials: map[string]interface{}{
					"protocols": protos,
				},
			}
		)

		It("returns expected error", func() {
			_, err := rabbitmq.URI(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid AMQP protocol credentials"))
		})
	})

	Context("when amqp uris field is missing", func() {
		var (
			protos = map[string]interface{}{
				"amqp": map[string]interface{}{
					"bar": false,
				},
			}

			svc = service.Service{
				Credentials: map[string]interface{}{
					"protocols": protos,
				},
			}
		)

		It("returns expected error", func() {
			_, err := rabbitmq.URI(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid AMQP URIs"))
		})
	})

	Context("when amqp uris array is empty", func() {
		var (
			protos = map[string]interface{}{
				"amqp": map[string]interface{}{
					"uris": []interface{}{},
				},
			}

			svc = service.Service{
				Credentials: map[string]interface{}{
					"protocols": protos,
				},
			}
		)

		It("returns expected error", func() {
			_, err := rabbitmq.URI(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid AMQP URIs"))
		})
	})

	Context("when amqp uri is an empty string", func() {
		var (
			protos = map[string]interface{}{
				"amqp": map[string]interface{}{
					"uris": []interface{}{""},
				},
			}

			svc = service.Service{
				Credentials: map[string]interface{}{
					"protocols": protos,
				},
			}
		)

		It("returns expected error", func() {
			_, err := rabbitmq.URI(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid AMQP URI"))
		})
	})
})
