package service

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/cfkit/env"
	"github.com/streadway/amqp"
)

var _ = Describe(".Rabbit", func() {
	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServices)
	})

	It("does not return an error", func() {
		_, err := Rabbit()
		Expect(err).ToNot(HaveOccurred())
	})

	It("correctly initializes the amqp URI", func() {
		rabbit, _ := Rabbit()
		Expect(rabbit.uri).To(Equal("amqp://username:password@127.0.0.1/instance"))
	})

	Context("when service is not found", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_SERVICES", "{}")
		})

		It("returns an error", func() {
			_, err := Rabbit()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})
})

var _ = Describe(".RabbitWithTag", func() {
	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServices)
	})

	It("does not return an error", func() {
		_, err := RabbitWithTag("my-rabbit-tag")
		Expect(err).ToNot(HaveOccurred())
	})

	It("correctly initializes the amqp URI", func() {
		rabbit, _ := RabbitWithTag("my-rabbit-tag")
		Expect(rabbit.URI()).To(Equal("amqp://my-rabbit"))
	})

	Context("when service is not found", func() {
		It("returns an error", func() {
			_, err := RabbitWithTag("unknown")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Context("when getting the URI fails", func() {
		var (
			origLift    = serviceLift
			expectedErr = errors.New("expected")
		)

		BeforeEach(func() {
			serviceLift = func(s env.Service) (*RabbitMQ, error) {
				return nil, expectedErr
			}
		})

		AfterEach(func() {
			serviceLift = origLift
		})

		It("returns the epected error", func() {
			_, err := RabbitWithTag("my-rabbit-tag")
			Expect(err).To(Equal(expectedErr))
		})
	})
})

var _ = Describe(".RabbitWithName", func() {
	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServices)
	})

	It("does not return an error", func() {
		_, err := RabbitWithName("my-rabbit-name")
		Expect(err).ToNot(HaveOccurred())
	})

	It("correctly initializes the amqp URI", func() {
		rabbit, _ := RabbitWithName("my-rabbit-name")
		Expect(rabbit.URI()).To(Equal("amqp://my-rabbit"))
	})

	Context("when service is not found", func() {
		It("returns an error", func() {
			_, err := RabbitWithName("unknown")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Context("when parsing the service credentials fails", func() {
		var (
			origLift    = serviceLift
			expectedErr = errors.New("expected")
		)

		BeforeEach(func() {
			serviceLift = func(s env.Service) (*RabbitMQ, error) {
				return nil, expectedErr
			}
		})

		AfterEach(func() {
			serviceLift = origLift
		})

		It("returns the epected error", func() {
			_, err := RabbitWithName("my-rabbit-name")
			Expect(err).To(Equal(expectedErr))
		})
	})
})

var _ = Describe(".RabbitFromService", func() {
	var svc env.Service

	BeforeEach(func() {
		svc = env.Service{
			Label: "prabbit",
			Credentials: map[string]interface{}{
				"uri": "amqp://uri",
			},
		}
	})

	It("does not return an error", func() {
		_, err := RabbitFromService(svc)
		Expect(err).ToNot(HaveOccurred())
	})

	It("uses the correct URI", func() {
		rabbit, _ := RabbitFromService(svc)
		Expect(rabbit.uri).To(Equal("amqp://uri"))
	})

	Context("when service credentials do NOT contain URI", func() {
		BeforeEach(func() {
			svc.Credentials = map[string]interface{}{}
		})

		It("returns expected error", func() {
			_, err := RabbitFromService(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid AMQP URI"))
		})
	})

	Context("when service URI is empty", func() {
		BeforeEach(func() {
			svc.Credentials = map[string]interface{}{
				"uri": "",
			}
		})

		It("returns expected error", func() {
			_, err := RabbitFromService(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid AMQP URI"))
		})
	})

	Context("when service URI does not specify amqp protocol", func() {
		BeforeEach(func() {
			svc.Credentials = map[string]interface{}{
				"uri": "http://foo.bar",
			}
		})

		It("returns expected error", func() {
			_, err := RabbitFromService(svc)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Invalid AMQP URI"))
		})
	})
})

var _ = Describe("RabbitMQ", func() {
	Describe(".URI", func() {
		It("returns the correct URI", func() {
			r := &RabbitMQ{"uri"}
			Expect(r.URI()).To(Equal("uri"))
		})
	})

	Describe(".Dial", func() {
		var (
			dialedURI string

			actualConn *amqp.Connection
			actualErr  error

			expectedConn *amqp.Connection
			expectedErr  error

			origDialer = amqpDialer

			testDialer = func(uri string) (*amqp.Connection, error) {
				dialedURI = uri
				return expectedConn, expectedErr
			}
		)

		BeforeEach(func() {
			expectedConn = new(amqp.Connection)
			expectedErr = nil
			amqpDialer = testDialer
		})

		AfterEach(func() {
			amqpDialer = origDialer
		})

		JustBeforeEach(func() {
			rabbit := &RabbitMQ{"amqp://rabbit.uri"}
			actualConn, actualErr = rabbit.Dial()
		})

		It("uses the correct amqp URI", func() {
			Expect(dialedURI).To(Equal("amqp://rabbit.uri"))
		})

		It("returns the expected connection", func() {
			Expect(actualConn).To(Equal(expectedConn))
		})

		It("returns no error", func() {
			Expect(actualErr).ToNot(HaveOccurred())
		})

		Context("when the dialer fails", func() {
			BeforeEach(func() {
				expectedErr = errors.New("some-error")
			})

			It("returns the error", func() {
				Expect(actualErr).To(Equal(expectedErr))
			})
		})
	})
})

var vcapServices = `
	{
		"p-rabbitmq": [
		 {
			"credentials": {
			 "dashboard_url": "https://dashboard.foo.bar/#/login/username/password",
			 "hostname": "127.0.0.1",
			 "hostnames": [
				"127.0.0.1",
				"127.0.0.2",
				"127.0.0.3",
				"127.0.0.4"
			 ],
			 "http_api_uri": "https://username:password@pivotal-rabbitmq.run.pez.pivotal.io/api/",
			 "http_api_uris": [
				"https://username:password@pivotal-rabbitmq.run.pez.pivotal.io/api/"
			 ],
			 "password": "password",
			 "protocols": {
				"amqp": {
				 "host": "127.0.0.1",
				 "hosts": [
					"127.0.0.1",
					"127.0.0.2",
					"127.0.0.3",
					"127.0.0.4"
				 ],
				 "password": "password",
				 "port": 5672,
				 "ssl": false,
				 "uri": "amqp://username:password@127.0.0.1:5672/instance",
				 "uris": [
					"amqp://username:password@127.0.0.1:5672/instance",
					"amqp://username:password@127.0.0.2:5672/instance",
					"amqp://username:password@127.0.0.3:5672/instance",
					"amqp://username:password@127.0.0.4:5672/instance"
				 ],
				 "username": "username",
				 "vhost": "instance"
				},
				"management": {
				 "host": "127.0.0.1",
				 "hosts": [
					"127.0.0.1",
					"127.0.0.2",
					"127.0.0.3",
					"127.0.0.4"
				 ],
				 "password": "password",
				 "path": "/api/",
				 "port": 15672,
				 "ssl": false,
				 "uri": "http://username:password@127.0.0.1:15672/api/",
				 "uris": [
					"http://username:password@127.0.0.1:15672/api/",
					"http://username:password@127.0.0.2:15672/api/",
					"http://username:password@127.0.0.3:15672/api/",
					"http://username:password@127.0.0.4:15672/api/"
				 ],
				 "username": "username"
				},
				"mqtt": {
				 "host": "127.0.0.1",
				 "hosts": [
					"127.0.0.1",
					"127.0.0.2",
					"127.0.0.3",
					"127.0.0.4"
				 ],
				 "password": "password",
				 "port": 1883,
				 "ssl": false,
				 "uri": "mqtt://instance%3Ausername:password@127.0.0.1:1883",
				 "uris": [
					"mqtt://instance%3Ausername:password@127.0.0.1:1883",
					"mqtt://instance%3Ausername:password@127.0.0.2:1883",
					"mqtt://instance%3Ausername:password@127.0.0.3:1883",
					"mqtt://instance%3Ausername:password@127.0.0.4:1883"
				 ],
				 "username": "instance:username"
				},
				"stomp": {
				 "host": "127.0.0.1",
				 "hosts": [
					"127.0.0.1",
					"127.0.0.2",
					"127.0.0.3",
					"127.0.0.4"
				 ],
				 "password": "password",
				 "port": 61613,
				 "ssl": false,
				 "uri": "stomp://username:password@127.0.0.1:61613",
				 "uris": [
					"stomp://username:password@127.0.0.1:61613",
					"stomp://username:password@127.0.0.2:61613",
					"stomp://username:password@127.0.0.3:61613",
					"stomp://username:password@127.0.0.4:61613"
				 ],
				 "username": "username",
				 "vhost": "instance"
				}
			 },
			 "ssl": false,
			 "uri": "amqp://username:password@127.0.0.1/instance",
			 "uris": [
				"amqp://username:password@127.0.0.1/instance",
				"amqp://username:password@127.0.0.2/instance",
				"amqp://username:password@127.0.0.3/instance",
				"amqp://username:password@127.0.0.4/instance"
			 ],
			 "username": "username",
			 "vhost": "instance"
			},
			"label": "p-rabbitmq",
			"name": "rabbit",
			"plan": "standard",
			"tags": [
			 "rabbitmq",
			 "messaging",
			 "message-queue",
			 "amqp",
			 "stomp",
			 "mqtt",
			 "pivotal"
			]
		 }
		],
		"cloudamqp-dev": [
		 {
			"credentials": {
				"uri": "amqp://my-rabbit"
			},
			"label": "cloudamqp-dev",
			"name": "my-rabbit-name",
			"syslog_drain_url": "",
			"tags": ["my-rabbit-tag"]
		 }
	  ]
	}
`
