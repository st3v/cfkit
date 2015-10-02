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

	defaultHeartbeatInterval = 30 * time.Second
)

var connProvider = fargoConn

func NewClient(uris []string, port int, timeout, pollInterval time.Duration) *Client {
	return &Client{
		conn:              connProvider(uris, port, timeout, pollInterval),
		uris:              uris,
		port:              port,
		timeout:           timeout,
		pollInterval:      pollInterval,
		heartbeatInterval: defaultHeartbeatInterval,
	}
}

type Client struct {
	conn              FargoConnection
	uris              []string
	port              int
	timeout           time.Duration
	pollInterval      time.Duration
	heartbeatInterval time.Duration
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

func (c *Client) Register(app env.App) error {
	instance := fargoInstance(app)
	err := c.conn.RegisterInstance(instance)
	if err != nil {
		return fmt.Errorf("Error registering app with Eureka: %s", err)
	}

	// assume all apps use the same renewal interval
	if instance.LeaseInfo.RenewalIntervalInSecs > 0 {
		c.heartbeatInterval = time.Duration(instance.LeaseInfo.RenewalIntervalInSecs) * time.Second
	}

	return nil
}

func (c *Client) Deregister(app env.App) error {
	err := c.conn.DeregisterInstance(fargoInstance(app))
	if err != nil {
		return fmt.Errorf("Error deregistering app with Eureka: %s", err)
	}
	return nil
}

func (c *Client) Heartbeat(app env.App) error {
	err := c.conn.HeartBeatInstance(fargoInstance(app))
	if err != nil {
		return fmt.Errorf("Error sending heartbeat for app to Eureka: %s", err)
	}
	return nil
}

func (c *Client) HeartbeatInterval() time.Duration {
	return c.heartbeatInterval
}

func (c *Client) Apps() (map[string][]string, error) {
	apps, err := c.conn.GetApps()
	if err != nil {
		return map[string][]string{}, fmt.Errorf("Error retrieving apps from Eureka: %s", err)
	}

	result := map[string][]string{}
	for _, app := range apps {
		uris := make([]string, len(app.Instances))

		for i, inst := range app.Instances {
			uris[i] = inst.HostName
		}

		result[app.Name] = uris
	}

	return result, nil
}

func (c *Client) App(name string) ([]string, error) {
	app, err := c.conn.GetApp(name)
	if err != nil {
		return []string{}, fmt.Errorf("Error retrieving app '%s' from Eureka: %s", name, err)
	}

	result := make([]string, len(app.Instances))
	for i, inst := range app.Instances {
		result[i] = inst.HostName
	}

	return result, nil
}

func (c *Client) URIs() []string {
	return c.uris
}

func (c *Client) Port() int {
	return c.port
}

func (c *Client) Timeout() time.Duration {
	return c.timeout
}

func (c *Client) PollInterval() time.Duration {
	return c.pollInterval
}
