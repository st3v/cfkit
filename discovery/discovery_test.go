package discovery

import (
	"errors"
	"io"
	"log"
	"os"
	"sync/atomic"
	"syscall"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/gbytes"
	"github.com/st3v/cfkit/discovery/fake"
	"github.com/st3v/cfkit/env"
)

var _ = Describe("discovery", func() {
	var (
		fakeClient         *fake.Client
		origClientProvider = clientProvider
		expectedApp        env.App
		exitCode           int32
	)

	BeforeEach(func() {
		os.Setenv("VCAP_APPLICATION", vcapApplication)

		// override os.Exit in tests
		atomic.StoreInt32(&exitCode, 0)
		exit = func(i int) {
			atomic.AddInt32(&exitCode, 1)
		}

		fakeClient = new(fake.Client)
		fakeClient.HeartbeatIntervalReturns(10 * time.Millisecond)

		clientProvider = func() (Client, error) {
			return fakeClient, nil
		}

		var err error
		expectedApp, err = env.Application()
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		Disable()
		clientProvider = origClientProvider
		os.Unsetenv("VCAP_APPLICATION")
	})

	Describe(".Enable", func() {
		It("registers the app with eureka", func() {
			Enable()
			Eventually(fakeClient.RegisterCallCount).Should(Equal(1))
			Expect(fakeClient.RegisterArgsForCall(0)).To(Equal(expectedApp))
		})

		It("sends regular heartbeats", func() {
			Enable()
			Eventually(fakeClient.HeartbeatCallCount).Should(BeNumerically(">=", 10))
			for i := 0; i < fakeClient.HeartbeatCallCount(); i++ {
				Expect(fakeClient.HeartbeatArgsForCall(i)).To(Equal(expectedApp))
			}
		})

		It("deregisters the app upon receiving SIGINT", func() {
			Enable()
			Consistently(fakeClient.DeregisterCallCount).Should(Equal(0))
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			Eventually(fakeClient.DeregisterCallCount).Should(Equal(1))
			Expect(atomic.LoadInt32(&exitCode)).ToNot(BeZero())
		})

		It("deregisters the app upon receiving SIGTERM", func() {
			Enable()
			Consistently(fakeClient.DeregisterCallCount).Should(Equal(0))
			syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
			Eventually(fakeClient.DeregisterCallCount).Should(Equal(1))
			Expect(atomic.LoadInt32(&exitCode)).ToNot(BeZero())
		})

		It("deregisters the app upon receiving SIGHUP", func() {
			Enable()
			Consistently(fakeClient.DeregisterCallCount).Should(Equal(0))
			syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
			Eventually(fakeClient.DeregisterCallCount).Should(Equal(1))
			Expect(atomic.LoadInt32(&exitCode)).ToNot(BeZero())
		})

		Context("when it is called multiple times without having been disabled", func() {
			var clientProviderCalled bool

			BeforeEach(func() {
				Enable()

				clientProviderCalled = false

				clientProvider = func() (Client, error) {
					clientProviderCalled = true
					return fakeClient, nil
				}
			})

			It("does not do anything", func() {
				Enable()
				Expect(clientProviderCalled).To(BeFalse())
			})
		})

		Context("when registering the app fails", func() {
			BeforeEach(func() {
				retryTimeout = 10 * time.Millisecond
				fakeClient.RegisterReturns(errors.New("some-error"))
			})

			It("keeps retrying", func() {
				Enable()
				Eventually(fakeClient.RegisterCallCount).Should(BeNumerically(">=", 10))
			})

			It("does not send heartbeats", func() {
				Enable()
				Consistently(fakeClient.HeartbeatCallCount).Should(Equal(0))
			})
		})

		Context("when sending the heartbeat fails", func() {
			BeforeEach(func() {
				retryTimeout = 10 * time.Millisecond
				fakeClient.HeartbeatReturns(errors.New("some-error"))
			})

			It("reregisters the app", func() {
				Enable()
				Eventually(fakeClient.RegisterCallCount).Should(BeNumerically(">=", 10))
			})

			It("keeps retrying", func() {
				Enable()
				Eventually(fakeClient.HeartbeatCallCount).Should(BeNumerically(">=", 10))
			})
		})

		Context("when getting the eureka client fails", func() {
			var (
				expectedErr = errors.New("some-error")
				buf         io.ReadWriter
			)

			BeforeEach(func() {
				buf = gbytes.NewBuffer()
				log.SetOutput(io.MultiWriter(buf, GinkgoWriter))

				clientProvider = func() (Client, error) {
					return fakeClient, expectedErr
				}
			})

			AfterEach(func() {
				log.SetOutput(GinkgoWriter)
			})

			It("exits with a non-zero exit code", func() {
				Enable()
				Expect(atomic.LoadInt32(&exitCode)).ToNot(BeZero())
			})

			It("logs the error", func() {
				Enable()
				Expect(buf).To(gbytes.Say(expectedErr.Error()))
			})
		})

		Context("when getting the app from env fails", func() {
			var (
				origAppProvider = appProvider
				expectedErr     = errors.New("some-error")
				buf             io.ReadWriter
			)

			BeforeEach(func() {
				buf = gbytes.NewBuffer()
				log.SetOutput(io.MultiWriter(buf, GinkgoWriter))

				appProvider = func() (env.App, error) {
					return env.App{}, expectedErr
				}
			})

			AfterEach(func() {
				appProvider = origAppProvider
				log.SetOutput(GinkgoWriter)
			})

			It("exits with a non-zero exit code", func() {
				Enable()
				Expect(atomic.LoadInt32(&exitCode)).ToNot(BeZero())
			})

			It("logs the error", func() {
				Enable()
				Expect(buf).To(gbytes.Say(expectedErr.Error()))
			})
		})
	})

	Describe(".Disable", func() {
		BeforeEach(func() {
			Enable()
		})

		It("deregisters the app", func() {
			Consistently(fakeClient.DeregisterCallCount).Should(Equal(0))
			Disable()
			Eventually(fakeClient.DeregisterCallCount).Should(Equal(1))
		})

		It("stops registering the app", func() {
			Disable()
			callCount := fakeClient.RegisterCallCount()
			Consistently(fakeClient.RegisterCallCount).Should(Equal(callCount))
		})

		It("stops sending heartbeats for the app", func() {
			Disable()
			callCount := fakeClient.HeartbeatCallCount()
			Consistently(fakeClient.HeartbeatCallCount).Should(Equal(callCount))
		})

		It("cancels the internal context", func() {
			cancelCalled := false
			cancel = func() {
				cancelCalled = true
			}
			Disable()
			Expect(cancelCalled).To(BeTrue())
		})

		It("sets cancel to nil", func() {
			Expect(cancel).ToNot(BeNil())
			Disable()
			Expect(cancel).To(BeNil())
		})

		Context("when it is called multiple times", func() {
			BeforeEach(func() {
				Disable()
				Eventually(fakeClient.DeregisterCallCount).Should(Equal(1))
			})

			It("does not do anything", func() {
				Disable()
				Consistently(fakeClient.DeregisterCallCount).Should(Equal(1))
			})
		})
	})
})

var _ = Describe("clientProvider", func() {
	Context("when a Eureka service exists", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_SERVICES", vcapServices)
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_SERVICES")
		})

		It("returns a client", func() {
			c, _ := clientProvider()
			Expect(c).ToNot(BeNil())
		})

		It("returns no error", func() {
			_, err := clientProvider()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when a Eureka service does not exist", func() {
		It("returns no error", func() {
			_, err := clientProvider()
			Expect(err).To(HaveOccurred())
		})
	})
})

var vcapApplication = `
{
  "application_id": "e16ad474-0e22-42d4-98c7-d41ed0eec123",
  "application_name": "cfkit",
  "application_uris": [
   "cfkit.cfapps.io"
  ],
  "application_version": "e53f75c2-3723-47dd-b988-67c296a998ca",
	"host": "0.0.0.0",
  "port": 63940,
  "limits": {
   "disk": 1024,
   "fds": 16384,
   "mem": 64
  },
  "name": "cfkit",
  "instance_id": "3fc7db2dfa534d3cb6094f17fe6e12f5",
  "instance_index": 99,
  "space_id": "cc35031c-b4af-4eea-9914-b25cc0db3888",
  "space_name": "development",
  "uris": [
   "cfkit.cfapps.io"
  ],
  "users": null,
  "version": "e53f75c2-3723-47dd-b988-67c296a998ca",
	"started_at_timestamp": 123456789,
	"state_timestamp": 987654321
 }
`

var vcapServices = `
{
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
   }
  ]
}`
