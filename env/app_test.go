package env_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/cfkit/env"
)

var _ = Describe(".Application", func() {
	Context("when all relevant env vars are set", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_APPLICATION", vcapApplication)
			os.Setenv("CF_INSTANCE_INDEX", "99")
			os.Setenv("CF_INSTANCE_IP", "1.2.3.4")
			os.Setenv("CF_INSTANCE_PORT", "12345")
			os.Setenv("CF_INSTANCE_ADDR", "1.2.3.4:12345")
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_APPLICATION")
			os.Unsetenv("CF_INSTANCE_INDEX")
			os.Unsetenv("CF_INSTANCE_ID")
			os.Unsetenv("CF_INSTANCE_PORT")
			os.Unsetenv("CF_INSTANCE_ADDR")
		})

		It("does not return an error", func() {
			_, err := env.Application()
			Expect(err).ToNot(HaveOccurred())
		})

		It("correctly initializes all values", func() {
			app, _ := env.Application()
			Expect(app.ID).To(Equal("e16ad474-0e22-42d4-98c7-d41ed0eec123"))
			Expect(app.Name).To(Equal("cfkit"))
			Expect(app.URIs).To(HaveLen(1))
			Expect(app.URIs[0]).To(Equal("cfkit.cfapps.io"))
			Expect(app.Host).To(Equal("0.0.0.0"))
			Expect(app.Port).To(Equal(63940))
			Expect(app.Addr).To(Equal("0.0.0.0:63940"))
			Expect(app.Version).To(Equal("e53f75c2-3723-47dd-b988-67c296a998ca"))
			Expect(app.Limits.Memory).To(Equal(64))
			Expect(app.Limits.Disk).To(Equal(1024))
			Expect(app.Limits.FileDescriptors).To(Equal(16384))
			Expect(app.Space.ID).To(Equal("cc35031c-b4af-4eea-9914-b25cc0db3888"))
			Expect(app.Space.Name).To(Equal("development"))
			Expect(app.StartTimestamp).To(Equal(123456789))
			Expect(app.StateTimestamp).To(Equal(987654321))
			Expect(app.Instance.Index).To(Equal(99))
			Expect(app.Instance.IP).To(Equal("1.2.3.4"))
			Expect(app.Instance.Port).To(Equal(12345))
			Expect(app.Instance.Addr).To(Equal("1.2.3.4:12345"))
		})
	})

	Context("when CF_INSTANCE_INDEX env var is not set", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_APPLICATION", vcapApplication)
			os.Unsetenv("CF_INSTANCE_INDEX")
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_APPLICATION")
		})

		It("does not return an error", func() {
			_, err := env.Application()
			Expect(err).ToNot(HaveOccurred())
		})

		It("sets instance index to zero", func() {
			app, _ := env.Application()
			Expect(app.Instance.Index).To(BeZero())
		})
	})

	Context("when CF_INSTANCE_IP env var is not set", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_APPLICATION", vcapApplication)
			os.Unsetenv("CF_INSTANCE_IP")
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_APPLICATION")
		})

		It("does not return an error", func() {
			_, err := env.Application()
			Expect(err).ToNot(HaveOccurred())
		})

		It("sets an empty instance IP", func() {
			app, _ := env.Application()
			Expect(app.Instance.IP).To(Equal(""))
		})
	})

	Context("when CF_INSTANCE_PORT env var is not set", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_APPLICATION", vcapApplication)
			os.Unsetenv("CF_INSTANCE_PORT")
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_APPLICATION")
		})

		It("does not return an error", func() {
			_, err := env.Application()
			Expect(err).ToNot(HaveOccurred())
		})

		It("sets instance port to zero", func() {
			app, _ := env.Application()
			Expect(app.Instance.Port).To(BeZero())
		})
	})

	Context("when CF_INSTANCE_ADDR env var is not set", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_APPLICATION", vcapApplication)
			os.Unsetenv("CF_INSTANCE_ADDR")
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_APPLICATION")
		})

		It("does not return an error", func() {
			_, err := env.Application()
			Expect(err).ToNot(HaveOccurred())
		})

		It("sets an empty instance Addr", func() {
			app, _ := env.Application()
			Expect(app.Instance.Addr).To(Equal(""))
		})
	})

	Context("when env var VCAP_APPLICATION is not set", func() {
		BeforeEach(func() {
			os.Unsetenv("VCAP_APPLICATION")
		})

		It("returns an error", func() {
			_, err := env.Application()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("VCAP_APPLICATION not set"))
		})
	})

	Context("when env var VCAP_APPLICATION is invalid", func() {
		BeforeEach(func() {
			os.Setenv("VCAP_APPLICATION", `{"application_id": 1234.5}`)
		})

		AfterEach(func() {
			os.Unsetenv("VCAP_APPLICATION")
		})

		It("returns an error", func() {
			_, err := env.Application()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Error parsing VCAP_APPLICATION"))
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
