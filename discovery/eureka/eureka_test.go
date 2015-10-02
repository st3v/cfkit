package eureka

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hudl/fargo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/st3v/cfkit/discovery/eureka/fake"
	"github.com/st3v/cfkit/env"
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

			conn, ok := c.conn.(*fargo.EurekaConnection)
			Expect(ok).To(BeTrue())

			Expect(conn.ServiceUrls).To(HaveLen(len(expectedURIs)))
			Expect(conn.ServiceUrls).To(ConsistOf(expectedURIs))
			Expect(conn.ServicePort).To(Equal(expectedPort))
			Expect(conn.Timeout).To(Equal(expectedTimeout))
			Expect(conn.PollInterval).To(Equal(expectedPollInterval))
		})
	})

	Describe("Client", func() {
		var (
			client *Client

			fakeConn         *fake.FargoConnection
			origConnProvider = fargoConn

			fargoApp  *fargo.Application
			fargoApps map[string]*fargo.Application

			app = env.App{
				ID:   "app-id",
				Name: "app-name",
				URIs: []string{"app-uri-1", "app-uri-2"},
				Host: "app-host",
				Port: 12345,
				Addr: "app-addr",
				Instance: env.AppInstance{
					ID:    "app-instance_id",
					Index: 99,
				},
			}

			assertInstance = func(i *fargo.Instance) {
				Expect(i.HostName).To(Equal(app.URI()))
				Expect(i.PortJ.Number).To(Equal("80"))
				Expect(i.PortJ.Enabled).To(Equal("true"))
				Expect(i.SecurePortJ.Number).To(Equal("443"))
				Expect(i.SecurePortJ.Enabled).To(Equal("true"))
				Expect(i.App).To(Equal(strings.ToUpper(app.Name)))
				Expect(i.IPAddr).To(Equal(app.Instance.Addr))
				Expect(i.VipAddress).To(Equal(app.URI()))
				Expect(i.Status).To(Equal(fargo.UP))
				Expect(i.Overriddenstatus).To(Equal(fargo.UNKNOWN))
				Expect(i.DataCenterInfo.Name).To(Equal(fargo.MyOwn))
				Expect(i.UniqueID).ToNot(BeNil())
				Expect(i.UniqueID(*i)).To(Equal(fmt.Sprintf("%s:%s", app.URI(), app.Instance.ID)))
			}
		)

		BeforeEach(func() {
			fakeConn = new(fake.FargoConnection)

			fargoApp = &fargo.Application{
				Name: app.Name,
				Instances: []*fargo.Instance{{
					UniqueID: func(fargo.Instance) string {
						return app.Instance.ID
					},
					HostName: app.URI(),
				}},
			}
			fakeConn.GetAppReturns(fargoApp, nil)

			fargoApps = map[string]*fargo.Application{
				fargoApp.Name: fargoApp,
			}
			fakeConn.GetAppsReturns(fargoApps, nil)

			connProvider = func(uris []string, port int, timeout, pollInterval time.Duration) FargoConnection {
				return fakeConn
			}

			client = NewClient(expectedURIs, expectedPort, expectedTimeout, expectedPollInterval)
		})

		AfterEach(func() {
			connProvider = origConnProvider
		})

		Describe(".Register", func() {
			It("calls conn.RegisterInstance with the correct instance", func() {
				client.Register(app)
				Expect(fakeConn.RegisterInstanceCallCount()).To(Equal(1))
				assertInstance(fakeConn.RegisterInstanceArgsForCall(0))
			})

			Context("when conn.RegisterInstance returns an error", func() {
				var expectedErr = errors.New("some-error")

				BeforeEach(func() {
					fakeConn.RegisterInstanceReturns(expectedErr)
				})

				It("returns the error", func() {
					err := client.Register(app)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Error registering app"))
					Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
				})
			})
		})

		Describe(".Deregister", func() {
			It("calls conn.DeregisterInstance with the correct instance", func() {
				client.Deregister(app)
				Expect(fakeConn.DeregisterInstanceCallCount()).To(Equal(1))
				assertInstance(fakeConn.DeregisterInstanceArgsForCall(0))
			})

			Context("when conn.RegisterInstance returns an error", func() {
				var expectedErr = errors.New("some-error")

				BeforeEach(func() {
					fakeConn.DeregisterInstanceReturns(expectedErr)
				})

				It("returns the error", func() {
					err := client.Deregister(app)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Error deregistering app"))
					Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
				})
			})
		})

		Describe(".Heartbeat", func() {
			It("calls conn.HeartBeatInstance with the correct instance", func() {
				client.Heartbeat(app)
				Expect(fakeConn.HeartBeatInstanceCallCount()).To(Equal(1))
				assertInstance(fakeConn.HeartBeatInstanceArgsForCall(0))
			})

			Context("when conn.HeartBeatInstance returns an error", func() {
				var expectedErr = errors.New("some-error")

				BeforeEach(func() {
					fakeConn.HeartBeatInstanceReturns(expectedErr)
				})

				It("returns the error", func() {
					err := client.Heartbeat(app)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Error sending heartbeat for app"))
					Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
				})
			})
		})

		Describe(".HeartbeatInterval", func() {
			It("returns the default interval", func() {
				Expect(client.HeartbeatInterval()).To(Equal(30 * time.Second))
			})

			Context("when the eureka server requires a different interval", func() {
				var expectedInterval = 123 * time.Second

				BeforeEach(func() {
					fakeConn.RegisterInstanceStub = func(i *fargo.Instance) error {
						i.LeaseInfo = fargo.LeaseInfo{
							RenewalIntervalInSecs: int32(expectedInterval.Seconds()),
						}
						return nil
					}
				})

				Context("and it is called prior to the first register", func() {
					It("returns the default interval", func() {
						Expect(client.HeartbeatInterval()).To(Equal(30 * time.Second))
					})
				})

				Context("and it is called after the first register", func() {
					BeforeEach(func() {
						Expect(client.Register(app)).To(Succeed())
					})

					It("returns the interval set by the server", func() {
						Expect(client.HeartbeatInterval()).To(Equal(expectedInterval))
					})
				})
			})
		})

		Describe(".App", func() {
			It("calls conn.GetApp with the correct app name", func() {
				client.App("foo")
				Expect(fakeConn.GetAppCallCount()).To(Equal(1))
				Expect(fakeConn.GetAppArgsForCall(0)).To(Equal("foo"))
			})

			It("returns the app retrieved from conn.GetApp", func() {
				uris, _ := client.App("foo")
				Expect(uris).To(Equal([]string{app.URI()}))
			})

			Context("when conn.GetApp returns an error", func() {
				var expectedErr = errors.New("some-error")

				BeforeEach(func() {
					fakeConn.GetAppReturns(nil, expectedErr)
				})

				It("returns the error", func() {
					_, err := client.App("foo")
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Error retrieving app 'foo'"))
					Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
				})
			})
		})

		Describe(".Apps", func() {
			It("calls conn.GetApps", func() {
				client.Apps()
				Expect(fakeConn.GetAppsCallCount()).To(Equal(1))
			})

			It("returns the apps retrieved from conn.GetApps", func() {
				apps, _ := client.Apps()
				Expect(apps).To(HaveLen(len(fargoApps)))
				for name, uris := range apps {
					Expect(name).To(Equal(app.Name))
					Expect(uris).To(Equal([]string{app.URI()}))
				}
			})

			Context("when conn.GetApps returns an error", func() {
				var expectedErr = errors.New("some-error")

				BeforeEach(func() {
					fakeConn.GetAppsReturns(nil, expectedErr)
				})

				It("returns the error", func() {
					_, err := client.Apps()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Error retrieving apps"))
					Expect(err.Error()).To(ContainSubstring(expectedErr.Error()))
				})
			})
		})

		Describe(".URIs", func() {
			It("returns the correct URIs", func() {
				Expect(client.URIs()).To(HaveLen(len(expectedURIs)))
				Expect(client.URIs()).To(ConsistOf(expectedURIs))
			})
		})

		Describe(".Port", func() {
			It("returns the correct port", func() {
				Expect(client.Port()).To(Equal(expectedPort))
			})
		})

		Describe(".Timeout", func() {
			It("returns the correct timeout", func() {
				Expect(client.Timeout()).To(Equal(expectedTimeout))
			})
		})

		Describe(".PollInterval", func() {
			It("returns the correct poll intervall", func() {
				Expect(client.PollInterval()).To(Equal(expectedPollInterval))
			})
		})
	})
})
