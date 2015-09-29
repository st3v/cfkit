package service

import (
	"os"
	"time"

	"github.com/hudl/fargo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/cfkit/env"
)

var _ = Describe(".Eureka", func() {
	var (
		origLift = eurekaLift
		lifted   env.Service
	)

	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServicesEureka)

		lifted = env.Service{}

		eurekaLift = func(svc env.Service) (*cfargo, error) {
			lifted = svc
			return origLift(svc)
		}
	})

	AfterEach(func() {
		eurekaLift = origLift
		os.Unsetenv("VCAP_SERVICES")
	})

	It("lifts the correct service", func() {
		Eureka()
		Expect(lifted.Name).To(Equal("eureka"))
		Expect(lifted.Label).To(Equal("user-provided"))
		Expect(lifted.Tags).To(HaveLen(0))
	})

	It("returns no error", func() {
		_, err := Eureka()
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns a EurekaClient with the correct URIs", func() {
		c, _ := Eureka()
		Expect(c.URIs()).To(HaveLen(3))
		Expect(c.URIs()).To(ContainElement("http://eureka-one-a.cfapps.io/eureka"))
		Expect(c.URIs()).To(ContainElement("http://eureka-one-b.cfapps.io/eureka"))
		Expect(c.URIs()).To(ContainElement("http://eureka-one-c.cfapps.io/eureka"))
	})
})

var _ = Describe(".EurekaWithName", func() {
	var (
		name     = "eureka-service"
		origLift = eurekaLift
		lifted   env.Service
	)

	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServicesEureka)

		lifted = env.Service{}

		eurekaLift = func(svc env.Service) (*cfargo, error) {
			lifted = svc
			return origLift(svc)
		}
	})

	AfterEach(func() {
		eurekaLift = origLift
		os.Unsetenv("VCAP_SERVICES")
	})

	It("lifts the correct service", func() {
		EurekaWithName(name)
		Expect(lifted.Name).To(Equal(name))
		Expect(lifted.Label).To(Equal("user-provided"))
		Expect(lifted.Tags).To(HaveLen(0))
	})

	It("returns no error", func() {
		_, err := EurekaWithName(name)
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns a EurekaClient with the correct URIs", func() {
		c, _ := EurekaWithName(name)
		Expect(c.URIs()).To(HaveLen(1))
		Expect(c.URIs()).To(ContainElement("http://eureka-two.cfapps.io/eureka"))
	})

	Context("when the service cannot be found", func() {
		It("returns a corresponding error", func() {
			_, err := EurekaWithName("unknown")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})
})

var _ = Describe(".EurekaWithTag", func() {
	var (
		tag      = "eureka-tag"
		origLift = eurekaLift
		lifted   env.Service
	)

	BeforeEach(func() {
		os.Setenv("VCAP_SERVICES", vcapServicesEureka)

		lifted = env.Service{}

		eurekaLift = func(svc env.Service) (*cfargo, error) {
			lifted = svc
			return origLift(svc)
		}
	})

	AfterEach(func() {
		eurekaLift = origLift
		os.Unsetenv("VCAP_SERVICES")
	})

	It("lifts the correct service", func() {
		EurekaWithTag(tag)
		Expect(lifted.Name).To(Equal("eureka-name"))
		Expect(lifted.Label).To(Equal("user-provided"))
		Expect(lifted.Tags).To(HaveLen(1))
		Expect(lifted.Tags).To(ContainElement(tag))
	})

	It("returns no error", func() {
		_, err := EurekaWithTag(tag)
		Expect(err).ToNot(HaveOccurred())
	})

	It("returns a EurekaClient with the correct URIs", func() {
		c, _ := EurekaWithTag(tag)
		Expect(c.URIs()).To(HaveLen(1))
		Expect(c.URIs()).To(ContainElement("http://eureka-three.cfapps.io/eureka"))
	})

	Context("when the service cannot be found", func() {
		It("returns a corresponding error", func() {
			_, err := EurekaWithTag("unknown")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})
})

var _ = Describe(".EurekaFromService", func() {
	Context("when the service credentials include a uri field", func() {
		var (
			svc env.Service
			uri string
		)

		JustBeforeEach(func() {
			svc = env.Service{
				Credentials: map[string]interface{}{
					"uri": uri,
				},
			}
		})

		Context("that does specify protocol and path", func() {
			BeforeEach(func() {
				uri = "https://host:1234/path"
			})

			It("does not return an error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly initializes the URIs for the EurekaClient", func() {
				c, _ := EurekaFromService(svc)
				Expect(c.URIs()).To(HaveLen(1))
				Expect(c.URIs()[0]).To(Equal("https://host:1234/path"))
			})
		})

		Context("that does NOT specify a protocol and path", func() {
			BeforeEach(func() {
				uri = "host"
			})

			It("does not return an error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns a EurekaClient with an augmented URI", func() {
				c, _ := EurekaFromService(svc)
				Expect(c.URIs()).To(HaveLen(1))
				Expect(c.URIs()[0]).To(Equal("http://host/eureka"))
			})
		})

		Context("that does have an empty path", func() {
			BeforeEach(func() {
				uri = "http://host/"
			})

			It("does not return an error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns a EurekaClient with an augmented URI", func() {
				c, _ := EurekaFromService(svc)
				Expect(c.URIs()).To(HaveLen(1))
				Expect(c.URIs()[0]).To(Equal("http://host/eureka"))
			})
		})

		Context("that is empty", func() {
			BeforeEach(func() {
				uri = ""
			})

			It("returns a corresponding error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Empty service URI"))
			})
		})

		Context("that is NOT a valid URI", func() {
			BeforeEach(func() {
				uri = "&*(^^%^@#$"
			})

			It("returns a corresponding error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("Error parsing service URI"))
			})
		})

		Context("that has an invalid type", func() {
			It("returns a corresponding error", func() {
				svc := env.Service{
					Credentials: map[string]interface{}{
						"uri": 1234,
					},
				}

				_, err := EurekaFromService(svc)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Missing or invalid service URIs"))
			})
		})
	})

	Context("when the service credentials include a uris field", func() {
		var (
			svc  env.Service
			uris []interface{}
		)

		JustBeforeEach(func() {
			svc = env.Service{
				Credentials: map[string]interface{}{
					"uris": uris,
				},
			}
		})

		Context("with multiple URIs that include protocol and path", func() {
			BeforeEach(func() {
				uris = []interface{}{
					"https://host1:1234/path1",
					"https://host2:1234/path2",
					"https://host3:1234/path3",
				}
			})

			It("does not return an error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).ToNot(HaveOccurred())
			})

			It("returns a correctly initialized EurekaClient", func() {
				c, _ := EurekaFromService(svc)

				Expect(c.URIs()).To(HaveLen(3))
				Expect(c.URIs()).To(ContainElement(uris[0]))
				Expect(c.URIs()).To(ContainElement(uris[1]))
				Expect(c.URIs()).To(ContainElement(uris[2]))
			})
		})

		Context("with multiple URIs that are missing protocol or path", func() {
			BeforeEach(func() {
				uris = []interface{}{
					"host1:1234/path1",
					"https://host2:1234",
					"host3",
					"host4/path4",
					"abc://host5/",
				}
			})

			It("does not return an error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).ToNot(HaveOccurred())
			})

			It("correctly augments all uris", func() {
				c, _ := EurekaFromService(svc)
				Expect(c.URIs()).To(HaveLen(5))
				Expect(c.URIs()).To(ContainElement("http://host1:1234/path1"))
				Expect(c.URIs()).To(ContainElement("https://host2:1234/eureka"))
				Expect(c.URIs()).To(ContainElement("http://host3/eureka"))
				Expect(c.URIs()).To(ContainElement("http://host4/path4"))
				Expect(c.URIs()).To(ContainElement("abc://host5/eureka"))
			})
		})

		Context("with multiple URIs, one of them is empty", func() {
			BeforeEach(func() {
				uris = []interface{}{
					"http://host1:1234/path1",
					"",
					"http://host3:1234/path3",
				}
			})

			It("returns a corresponding error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Empty service URI"))
			})
		})

		Context("with multiple URIs, one of them has an invalid format", func() {
			BeforeEach(func() {
				uris = []interface{}{
					"http://host1:1234/path1",
					"://",
					"http://host3:1234/path3",
				}
			})

			It("returns a corresponding error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(HavePrefix("Error parsing service URI"))
			})
		})

		Context("with multiple URIs, one of them has an invalid type", func() {
			BeforeEach(func() {
				uris = []interface{}{
					"http://host1:1234/path1",
					12345,
					"http://host3:1234/path3",
				}
			})

			It("returns a corresponding error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Missing or invalid service URIs"))
			})
		})

		Context("with no URIs", func() {
			BeforeEach(func() {
				uris = []interface{}{}
			})

			It("returns a corresponding error", func() {
				_, err := EurekaFromService(svc)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Missing or invalid service URIs"))
			})
		})

		Context("with an invalid type", func() {
			It("returns a corresponding error", func() {
				svc := env.Service{
					Credentials: map[string]interface{}{
						"uris": 1234,
					},
				}

				_, err := EurekaFromService(svc)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Missing or invalid service URIs"))
			})
		})
	})

	Context("when the service does not have any optional properties", func() {
		var svc = env.Service{
			Credentials: map[string]interface{}{
				"uri": "http://my-host/eureka",
			},
		}

		It("uses the default port for the underlying connection", func() {
			c, _ := EurekaFromService(svc)
			Expect(c.conn.ServicePort).To(Equal(80))
		})

		It("uses the default timeout for the underlying connection", func() {
			c, _ := EurekaFromService(svc)
			Expect(c.conn.Timeout).To(Equal(10 * time.Second))
		})

		It("uses the default poll interval for the underlying connection", func() {
			c, _ := EurekaFromService(svc)
			Expect(c.conn.PollInterval).To(Equal(30 * time.Second))
		})
	})

	Context("when the service specifies optional properties", func() {
		var svc = env.Service{
			Credentials: map[string]interface{}{
				"uri":           "http://my-host/eureka",
				"port":          12345,
				"timeout":       67,
				"poll_interval": 89,
			},
		}

		It("uses the specified port property for the underlying connection", func() {
			c, _ := EurekaFromService(svc)
			Expect(c.conn.ServicePort).To(Equal(12345))
		})

		It("uses the specified timeout for the underlying connection", func() {
			c, _ := EurekaFromService(svc)
			Expect(c.conn.Timeout).To(Equal(67 * time.Second))
		})

		It("uses the specified poll_interval for the underlying connection", func() {
			c, _ := EurekaFromService(svc)
			Expect(c.conn.PollInterval).To(Equal(89 * time.Second))
		})
	})
})

var _ = Describe("cfargo", func() {
	Describe(".URIs", func() {
		var client *cfargo

		BeforeEach(func() {
			client = &cfargo{&fargo.EurekaConnection{
				ServiceUrls: []string{"uri1", "uri2", "uri3"},
			}}
		})

		It("returns the service URIs of the underlying connection", func() {
			Expect(client.URIs()).To(HaveLen(3))
			Expect(client.URIs()).To(ContainElement("uri1"))
			Expect(client.URIs()).To(ContainElement("uri2"))
			Expect(client.URIs()).To(ContainElement("uri3"))
		})

		Context("when the underlying connection is nil", func() {
			BeforeEach(func() {
				client = new(cfargo)
			})

			It("returns an empty slice", func() {
				Expect(client.URIs()).To(HaveLen(0))
			})
		})
	})
})

var vcapServicesEureka = `{
	"user-provided": [
	 {
		"credentials": {
		 "uris": [
			"eureka-one-a.cfapps.io",
			"eureka-one-b.cfapps.io",
			"eureka-one-c.cfapps.io"
		 ]
		},
		"label": "user-provided",
		"name": "eureka",
		"syslog_drain_url": "",
		"tags": []
	 },
	 {
		"credentials": {
		 "uri": "eureka-two.cfapps.io"
		},
		"label": "user-provided",
		"name": "eureka-service",
		"syslog_drain_url": "",
		"tags": []
	 },
	 {
		"credentials": {
		 "uri": "eureka-three.cfapps.io"
		},
		"label": "user-provided",
		"name": "eureka-name",
		"syslog_drain_url": "",
		"tags": [
		 "eureka-tag"
		]
	 }
	]
	}`
