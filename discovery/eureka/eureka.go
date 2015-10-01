package eureka

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hudl/fargo"
	"github.com/st3v/cfkit/env"
)

const (
	appStdPort       = 80
	appStdSecurePort = 443
)

var connProvider = fargoConn

func NewClient(uris []string, port int, timeout, pollInterval time.Duration) *client {
	return &client{
		conn:         connProvider(uris, port, timeout, pollInterval),
		uris:         uris,
		port:         port,
		timeout:      timeout,
		pollInterval: pollInterval,
	}
}

type Client interface {
	Register(app env.App) error
	Deregister(app env.App) error
	Heartbeat(app env.App) error
	Apps() ([]Application, error)
	App(name string) (Application, error)
	URIs() []string
	Port() int
	Timeout() time.Duration
	PollInterval() time.Duration
}

type Application struct {
	Name      string
	Instances []Instance
}

type Instance struct {
	ID  string
	URI string
}

type FargoConnection interface {
	RegisterInstance(*fargo.Instance) error
	DeregisterInstance(*fargo.Instance) error
	HeartBeatInstance(*fargo.Instance) error
	GetApp(string) (*fargo.Application, error)
	GetApps() (map[string]*fargo.Application, error)
}

func fargoConn(uris []string, port int, timeout, pollInterval time.Duration) FargoConnection {
	return &fargo.EurekaConnection{
		ServiceUrls:  uris,
		ServicePort:  port,
		PollInterval: pollInterval,
		Timeout:      timeout,
		UseJson:      true,
	}
}

type client struct {
	conn         FargoConnection
	uris         []string
	port         int
	timeout      time.Duration
	pollInterval time.Duration
}

func fargoInstance(app env.App) *fargo.Instance {
	return &fargo.Instance{
		HostName:         app.URI(),
		PortJ:            fargo.Port{strconv.Itoa(appStdPort), "true"},
		SecurePortJ:      fargo.Port{strconv.Itoa(appStdSecurePort), "true"},
		App:              strings.ToUpper(app.Name),
		IPAddr:           app.Instance.Addr,
		VipAddress:       app.URI(),
		Status:           fargo.UP,
		Overriddenstatus: fargo.UNKNOWN,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		Metadata:         fargo.InstanceMetadata{Raw: []byte(fmt.Sprintf(`{"instanceId": "%s"}`, app.Instance.ID))},
		UniqueID: func(i fargo.Instance) string {
			return fmt.Sprintf("%s:%s", app.URI(), app.Instance.ID)
		},
	}
}

func (c *client) Register(app env.App) error {
	err := c.conn.RegisterInstance(fargoInstance(app))
	if err != nil {
		return fmt.Errorf("Error registering app with Eureka: %s", err)
	}
	return nil
}

func (c *client) Deregister(app env.App) error {
	err := c.conn.DeregisterInstance(fargoInstance(app))
	if err != nil {
		return fmt.Errorf("Error deregistering app with Eureka: %s", err)
	}
	return nil
}

func (c *client) Heartbeat(app env.App) error {
	err := c.conn.HeartBeatInstance(fargoInstance(app))
	if err != nil {
		return fmt.Errorf("Error sending heartbeat for app to Eureka: %s", err)
	}
	return nil
}

func (c *client) Apps() ([]Application, error) {
	apps, err := c.conn.GetApps()
	if err != nil {
		return []Application{}, fmt.Errorf("Error retrieving apps from Eureka: %s", err)
	}

	result := []Application{}
	for _, app := range apps {
		instances := make([]Instance, len(app.Instances))

		for i, inst := range app.Instances {
			instances[i] = Instance{
				ID:  inst.Id(),
				URI: inst.VipAddress,
			}
		}

		result = append(result, Application{
			Name:      app.Name,
			Instances: instances,
		})
	}

	return result, nil
}

func (c *client) App(name string) (Application, error) {
	app, err := c.conn.GetApp(name)
	if err != nil {
		return Application{}, fmt.Errorf("Error retrieving app '%s' from Eureka: %s", name, err)
	}

	instances := make([]Instance, len(app.Instances))
	for i, inst := range app.Instances {
		instances[i] = Instance{
			ID:  inst.Id(),
			URI: inst.VipAddress,
		}
	}

	return Application{
		Name:      app.Name,
		Instances: instances,
	}, nil
}

func (c *client) URIs() []string {
	return c.uris
}

func (c *client) Port() int {
	return c.port
}

func (c *client) Timeout() time.Duration {
	return c.timeout
}

func (c *client) PollInterval() time.Duration {
	return c.pollInterval
}
