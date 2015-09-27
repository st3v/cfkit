package service_test

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cf "github.com/st3v/cfkit/integration_test/cfhelper"
)

var (
	api      = cf.API()
	username = cf.Username()
	password = cf.Password()
	org      = cf.Org()
	space    = cf.RandomSpaceName()
)

func TestService(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Integration Suite")
}

var _ = BeforeSuite(func() {
	Expect(cf.Login(api, username, password)).To(Succeed())
	Expect(cf.TargetOrg(org)).To(Succeed())
	Expect(cf.CreateSpace(space)).To(Succeed())
	Expect(cf.TargetSpace(space)).To(Succeed())
})

var _ = AfterSuite(func() {
	Expect(cf.DeleteSpace(space)).To(Succeed())
})
