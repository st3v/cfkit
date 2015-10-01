package eureka

import (
	"time"

	"github.com/hudl/fargo"
)

const (
	appStdPort       = 80
	appStdSecurePort = 443
)

func NewClient(uris []string, port int, timeout, pollInterval time.Duration) *fargoClient {
	return &fargoClient{fargoConn(uris, port, timeout, pollInterval)}
}

type Client interface {
	URIs() []string
	Port() int
	Timeout() time.Duration
	PollInterval() time.Duration
}

func fargoConn(uris []string, port int, timeout, pollInterval time.Duration) *fargo.EurekaConnection {
	return &fargo.EurekaConnection{
		ServiceUrls:  uris,
		ServicePort:  port,
		PollInterval: pollInterval,
		Timeout:      timeout,
		UseJson:      true,
	}
}

type fargoClient struct {
	conn *fargo.EurekaConnection
}

func (c *fargoClient) URIs() []string {
	if c.conn == nil {
		return []string{}
	}
	return c.conn.ServiceUrls
}

func (c *fargoClient) Port() int {
	if c.conn == nil {
		return 0
	}
	return c.conn.ServicePort
}

func (c *fargoClient) Timeout() time.Duration {
	if c.conn == nil {
		return 0
	}
	return c.conn.Timeout
}

func (c *fargoClient) PollInterval() time.Duration {
	if c.conn == nil {
		return 0
	}
	return c.conn.PollInterval
}
