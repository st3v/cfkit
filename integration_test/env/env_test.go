package env_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	cf "github.com/st3v/cfkit/integration_test/cfhelper"
)

func TestEnv(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Env Integration Suite")
}

var _ = Describe("Env", func() {
	var (
		api            = cf.API()
		username       = cf.Username()
		password       = cf.Password()
		org            = cf.Org()
		space          = cf.RandomSpaceName()
		app            = cf.RandomAppName()
		route          = fmt.Sprintf("http://%s.%s", app, cf.Domain())
		serviceID      = cf.RandomServiceID()
		servicePayload = map[string]interface{}{"uri": "some-uri"}
	)

	BeforeSuite(func() {
		Expect(cf.Login(api, username, password)).To(Succeed())
		Expect(cf.TargetOrg(org)).To(Succeed())

		Expect(cf.CreateSpace(space)).To(Succeed())
		Expect(cf.TargetSpace(space)).To(Succeed())

		manifestPath := filepath.Join(".", "env-test-app", "manifest.yml")
		Expect(cf.PushAppManifest(app, manifestPath)).To(Succeed())

		Expect(cf.CreateUserService(serviceID, servicePayload)).To(Succeed())
		Expect(cf.BindService(app, serviceID)).To(Succeed())

		Expect(cf.StartApp(app)).To(Succeed())
	})

	It("retrieves application properties from env", func() {
		resp, err := http.Get(fmt.Sprintf("%s/app", route))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var payload appPayload
		err = json.Unmarshal(body, &payload)
		Expect(err).ToNot(HaveOccurred())

		Expect(payload.Name).To(Equal(app))
		Expect(payload.URIs).To(HaveLen(1))
		Expect(payload.URIs[0]).To(Equal(fmt.Sprintf("%s.%s", app, cf.Domain())))
		Expect(payload.Version).ToNot(Equal(""))
		Expect(payload.Host).ToNot(Equal(""))
		Expect(payload.Port).To(BeNumerically(">", 0))
		Expect(payload.Addr).ToNot(Equal(""))
		Expect(payload.Limits.Memory).To(Equal(64))
		Expect(payload.Limits.Disk).To(BeNumerically(">", 0))
		Expect(payload.Limits.FileDescriptors).To(BeNumerically(">", 0))
		Expect(payload.StartTimestamp).To(BeNumerically(">", 0))
		Expect(payload.StateTimestamp).To(BeNumerically(">", 0))
		Expect(payload.Space.ID).ToNot(Equal(""))
		Expect(payload.Space.Name).To(Equal(space))
		Expect(payload.Instance.Index).To(BeNumerically(">=", 0))
		Expect(payload.Instance.IP).ToNot(Equal(""))
		Expect(payload.Instance.Port).To(BeNumerically(">", 0))
		Expect(payload.Instance.Addr).ToNot(Equal(""))
	})

	It("retrieves services from env", func() {
		resp, err := http.Get(fmt.Sprintf("%s/service/%s", route, serviceID))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var payload svcPayload
		err = json.Unmarshal(body, &payload)
		Expect(err).ToNot(HaveOccurred())

		Expect(payload.Label).To(Equal("user-provided"))
		Expect(payload.Name).To(Equal(serviceID))
		Expect(payload.Tags).To(BeEmpty())
		Expect(payload.Credentials.URI).To(Equal("some-uri"))
	})

	AfterSuite(func() {
		Expect(cf.DeleteApp(app)).To(Succeed())
		Expect(cf.DeleteService(serviceID)).To(Succeed())
		Expect(cf.DeleteSpace(space)).To(Succeed())
	})
})

type appPayload struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	URIs           []string `json:"uris"`
	Version        string   `json:"version"`
	Host           string   `json:"host"`
	Port           int      `json:"port"`
	Addr           string   `json:"addr"`
	StartTimestamp int      `json:"started_at_timestamp"`
	StateTimestamp int      `json:"state_timestamp"`
	Limits         struct {
		Memory          int `json:"mem"`
		Disk            int `json:"disk"`
		FileDescriptors int `json:"fds"`
	} `json:"limits"`
	Instance struct {
		Index int    `json:"index"`
		IP    string `json:"ip"`
		Port  int    `json:"port"`
		Addr  string `json:"addr"`
	} `json:"instance"`
	Space struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"space"`
}

type svcPayload struct {
	Name        string   `json:"name"`
	Label       string   `json:"label"`
	Tags        []string `json:"tags"`
	Credentials struct {
		URI string `json:"uri"`
	} `json:"credentials"`
}
