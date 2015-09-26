package rabbitmq_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cf "github.com/st3v/cfkit/integration_test/cfhelper"
)

var _ = Describe("RabbitMQ Service", func() {
	var (
		api         = cf.API()
		username    = cf.Username()
		password    = cf.Password()
		org         = cf.Org()
		space       = cf.RandomSpaceName()
		app         = cf.RandomAppName()
		serviceID   = cf.RandomServiceID()
		serviceName = cf.RabbitServiceName()
		servicePlan = cf.RabbitServicePlan()
		route       = fmt.Sprintf("http://%s.%s", app, cf.Domain())
	)

	BeforeSuite(func() {
		Expect(cf.Login(api, username, password)).To(Succeed())

		Expect(cf.TargetOrg(org)).To(Succeed())

		Expect(cf.CreateSpace(space)).To(Succeed())
		Expect(cf.TargetSpace(space)).To(Succeed())

		manifestPath := filepath.Join(".", "testapp", "manifest.yml")
		Expect(cf.PushAppManifest(app, manifestPath)).To(Succeed())

		Expect(cf.CreateService(
			serviceName,
			servicePlan,
			serviceID,
		)).To(Succeed())

		Expect(cf.BindService(app, serviceID)).To(Succeed())

		Expect(cf.StartApp(app)).To(Succeed())
	})

	It("posts and receives messages to and from a queue", func() {
		message := "Hello Rabbit!"

		resp, err := http.Post(route, "text/plain", strings.NewReader(message))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		fmt.Fprintln(GinkgoWriter, "Message posted")

		resp, err = http.Get(route)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(body)).To(Equal(message))

		fmt.Fprintln(GinkgoWriter, "Message received")
	})

	AfterSuite(func() {
		Expect(cf.DeleteApp(app)).To(Succeed())
		Expect(cf.DeleteService(serviceID)).To(Succeed())
		Expect(cf.DeleteSpace(space)).To(Succeed())
	})
})
