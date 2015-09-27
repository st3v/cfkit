package service_test

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
		app         = cf.RandomAppName()
		route       = fmt.Sprintf("http://%s.%s", app, cf.Domain())
		serviceID   = cf.RandomServiceID()
		serviceName = cf.RabbitServiceName()
		servicePlan = cf.RabbitServicePlan()
	)

	BeforeEach(func() {
		manifestPath := filepath.Join(".", "rabbit-test-app", "manifest.yml")
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

	AfterEach(func() {
		Expect(cf.DeleteApp(app)).To(Succeed())
		Expect(cf.DeleteService(serviceID)).To(Succeed())
	})
})
