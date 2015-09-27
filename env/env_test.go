package env_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/st3v/cfkit/env"
)

var _ = Describe(".Addr", func() {
	Context("when env var PORT is set", func() {
		BeforeEach(func() {
			os.Setenv("PORT", "1234")
		})

		AfterEach(func() {
			os.Unsetenv("PORT")
		})

		It("returns the expected address", func() {
			Expect(env.Addr()).To(Equal(":1234"))
		})
	})

	Context("when env var PORT is not set", func() {
		BeforeEach(func() {
			os.Unsetenv("PORT")
		})

		It("returns the expected address", func() {
			Expect(env.Addr()).To(Equal(":"))
		})
	})
})
