package rabbitmq

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/cfkit/service"
	"github.com/streadway/amqp"
)

var _ = Describe("rabbitmq", func() {
	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServices)
	})

	Describe(".Service", func() {
		It("does not return an error", func() {
			_, err := Service()
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly initializes the amqp URI", func() {
			rabbit, _ := Service()
			Expect(rabbit.uri).To(Equal("amqp://username:password@127.0.0.1:5672/instance"))
		})

		Context("when service is not found", func() {
			BeforeEach(func() {
				os.Setenv("VCAP_SERVICES", "{}")
			})

			It("returns an error", func() {
				_, err := Service()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})
	})

	Describe(".ServiceWithTag", func() {
		It("does not return an error", func() {
			_, err := ServiceWithTag("my-rabbit-tag")
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly initializes the amqp URI", func() {
			rabbit, _ := ServiceWithTag("my-rabbit-tag")
			Expect(rabbit.uri).To(Equal("amqp://my-rabbit"))
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				_, err := ServiceWithTag("unknown")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})

		Context("when getting the URI fails", func() {
			var (
				origURI     = URI
				expectedErr = errors.New("expected")
			)

			BeforeEach(func() {
				URI = func(s service.Service) (string, error) {
					return "", expectedErr
				}
			})

			AfterEach(func() {
				URI = origURI
			})

			It("returns the epected error", func() {
				_, err := ServiceWithTag("my-rabbit-tag")
				Expect(err).To(Equal(expectedErr))
			})
		})
	})

	Describe(".ServiceWithName", func() {
		It("does not return an error", func() {
			_, err := ServiceWithName("my-rabbit-name")
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly initializes the amqp URI", func() {
			rabbit, _ := ServiceWithName("my-rabbit-name")
			Expect(rabbit.uri).To(Equal("amqp://my-rabbit"))
		})

		Context("when service is not found", func() {
			It("returns an error", func() {
				_, err := ServiceWithName("unknown")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})

		Context("when getting the URI fails", func() {
			var (
				origURI     = URI
				expectedErr = errors.New("expected")
			)

			BeforeEach(func() {
				URI = func(s service.Service) (string, error) {
					return "", expectedErr
				}
			})

			AfterEach(func() {
				URI = origURI
			})

			It("returns the epected error", func() {
				_, err := ServiceWithName("my-rabbit-name")
				Expect(err).To(Equal(expectedErr))
			})
		})
	})

	Describe(".Dial", func() {
		var (
			dialedURI string

			actualConn *amqp.Connection
			actualErr  error

			expectedConn *amqp.Connection
			expectedErr  error

			origDialer func(string) (*amqp.Connection, error)

			testDialer = func(uri string) (*amqp.Connection, error) {
				dialedURI = uri
				return expectedConn, expectedErr
			}
		)

		BeforeEach(func() {
			expectedConn = new(amqp.Connection)
			expectedErr = nil

			origDialer = dialer
			dialer = testDialer
		})

		AfterEach(func() {
			dialer = origDialer
		})

		JustBeforeEach(func() {
			rabbit, err := Service()
			Expect(err).ToNot(HaveOccurred())

			actualConn, actualErr = rabbit.Dial()
		})

		It("uses the correct amqp URI", func() {
			Expect(dialedURI).To(Equal("amqp://username:password@127.0.0.1:5672/instance"))
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
		"user-provided": [
		 {
			"credentials": {
			 "protocols": {
				"amqp": {
					"uris": [
						"amqp://my-rabbit"
					]
				}
			 } 
			},
			"label": "user-provided",
			"name": "my-rabbit-name",
			"syslog_drain_url": "",
			"tags": ["my-rabbit-tag"]
		 }
	  ]
	}
`
