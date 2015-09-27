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
		Expect(rabbit.uri).To(Equal("amqp://username:password@127.0.0.1:5672/instance"))
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
			origFn      = rabbitFromService
			expectedErr = errors.New("expected")
		)

		BeforeEach(func() {
			rabbitFromService = func(s env.Service) (*RabbitMQ, error) {
				return nil, expectedErr
			}
		})

		AfterEach(func() {
			rabbitFromService = origFn
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
			origFn      = rabbitFromService
			expectedErr = errors.New("expected")
		)

		BeforeEach(func() {
			rabbitFromService = func(env.Service) (*RabbitMQ, error) {
				return nil, expectedErr
			}
		})

		AfterEach(func() {
			rabbitFromService = origFn
		})

		It("returns the epected error", func() {
			_, err := RabbitWithName("my-rabbit-name")
			Expect(err).To(Equal(expectedErr))
		})
	})
})

var _ = Describe("rabbitmq", func() {
	Describe("rabbit", func() {
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

	Describe(".rabbitFromService", func() {
		var svc env.Service

		Context("when the service is labeled p-rabbitmq", func() {
			BeforeEach(func() {
				protos := map[string]interface{}{
					"amqp": map[string]interface{}{
						"uri": "amqp://uri",
					},
				}

				svc = env.Service{
					Label: "p-rabbitmq",
					Credentials: map[string]interface{}{
						"protocols": protos,
					},
				}
			})

			It("does not return an error", func() {
				_, err := rabbitFromService(svc)
				Expect(err).ToNot(HaveOccurred())
			})

			It("uses the correct URI", func() {
				rabbit, _ := rabbitFromService(svc)
				Expect(rabbit.uri).To(Equal("amqp://uri"))
			})

			Context("when service protocols field is missing", func() {
				BeforeEach(func() {
					svc.Credentials = map[string]interface{}{}
				})

				It("returns expected error", func() {
					_, err := rabbitFromService(svc)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid service credentials"))
				})
			})

			Context("when protocols field is of an invalid type", func() {
				BeforeEach(func() {
					svc.Credentials = map[string]interface{}{
						"protocols": false,
					}
				})

				It("returns expected error", func() {
					_, err := rabbitFromService(svc)
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

					svc.Credentials = map[string]interface{}{
						"protocols": protos,
					}
				})

				It("returns expected error", func() {
					_, err := rabbitFromService(svc)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid AMQP protocol credentials"))
				})
			})

			Context("when the amqp field is of an invalid type", func() {
				BeforeEach(func() {
					protos := map[string]interface{}{
						"amqp": false,
					}

					svc.Credentials = map[string]interface{}{
						"protocols": protos,
					}
				})

				It("returns expected error", func() {
					_, err := rabbitFromService(svc)
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

					svc.Credentials = map[string]interface{}{
						"protocols": protos,
					}
				})

				It("returns expected error", func() {
					_, err := rabbitFromService(svc)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid AMQP URI"))
				})
			})

			Context("when amqp uris array is empty", func() {
				BeforeEach(func() {
					protos := map[string]interface{}{
						"amqp": map[string]interface{}{
							"uris": []interface{}{},
						},
					}

					svc.Credentials = map[string]interface{}{
						"protocols": protos,
					}
				})

				It("returns expected error", func() {
					_, err := rabbitFromService(svc)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid AMQP URI"))
				})
			})

			Context("when amqp uri is an empty string", func() {
				BeforeEach(func() {
					protos := map[string]interface{}{
						"amqp": map[string]interface{}{
							"uris": []interface{}{""},
						},
					}

					svc.Credentials = map[string]interface{}{
						"protocols": protos,
					}
				})

				It("returns expected error", func() {
					_, err := rabbitFromService(svc)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid AMQP URI"))
				})
			})
		})

		Context("when the service is NOT labeled p-rabbitmq", func() {
			BeforeEach(func() {
				svc = env.Service{
					Label: "cloudamqp",
					Credentials: map[string]interface{}{
						"uri": "amqp://uri",
					},
				}
			})

			It("does not return an error", func() {
				_, err := rabbitFromService(svc)
				Expect(err).ToNot(HaveOccurred())
			})

			It("uses the correct URI", func() {
				rabbit, _ := rabbitFromService(svc)
				Expect(rabbit.uri).To(Equal("amqp://uri"))
			})

			Context("when service credentials do NOT contain URI", func() {
				BeforeEach(func() {
					svc.Credentials = map[string]interface{}{}
				})

				It("returns expected error", func() {
					_, err := rabbitFromService(svc)
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
					_, err := rabbitFromService(svc)
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
					_, err := rabbitFromService(svc)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Invalid AMQP URI"))
				})
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
