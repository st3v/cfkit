package discovery

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/context"

	"github.com/st3v/cfkit/env"
	"github.com/st3v/cfkit/service"
)

var (
	retryTimeout                          = 10 * time.Second
	clientProvider func() (Client, error) = eurekaProvider
	appProvider                           = env.Application
	exit                                  = os.Exit
	cancel         context.CancelFunc
)

func eurekaProvider() (Client, error) {
	return service.Eureka()
}

type Client interface {
	Register(app env.App) error
	Deregister(app env.App) error
	Heartbeat(app env.App) error
	HeartbeatInterval() time.Duration
	Apps() (map[string][]string, error)
	App(name string) ([]string, error)
}

type Application struct {
	Name      string
	Instances []Instance
}

type Instance struct {
	ID  string
	URI string
}

func Disable() {
	if cancel != nil {
		cancel()
		cancel = nil
	}
}

func Enable() {
	if cancel != nil {
		return
	}

	client, err := clientProvider()
	if err != nil {
		log.Printf("Error getting discovery service: %s\n", err)
		exit(1)
	}

	app, err := appProvider()
	if err != nil {
		log.Printf("Error getting app from env: %s\n", err)
		exit(1)
	}

	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	deregisterOnShutdown(client, app, ctx, cancel, exit)
	registerAndKeepAlive(client, app, ctx)
}

func registerAndKeepAlive(client Client, app env.App, ctx context.Context) {
	go func(retryTimeout time.Duration) {
		// initial interval
		interval := 10 * time.Millisecond
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(interval):
				// subsequent interval
				interval = retryTimeout

				if err := client.Register(app); err != nil {
					log.Println(err.Error())
					continue
				}
				keepAlive(client, app, retryTimeout, ctx)
			}
		}
	}(retryTimeout)
}

func keepAlive(client Client, app env.App, retryTimeout time.Duration, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(client.HeartbeatInterval()):
			if err := client.Heartbeat(app); err != nil {
				log.Println(err.Error())
				time.Sleep(retryTimeout)
				return
			}
		}
	}
}

func deregisterOnShutdown(client Client, app env.App, ctx context.Context, cancel context.CancelFunc, exit func(int)) {
	sigChan := make(chan os.Signal, 1)

	signal.Reset(syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)

	go func() {
		select {
		case <-sigChan:
			cancel()
			client.Deregister(app)
		case <-ctx.Done():
			client.Deregister(app)
		}
		exit(1)
	}()
}
