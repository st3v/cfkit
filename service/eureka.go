package service

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/st3v/cfkit/discovery/eureka"
	"github.com/st3v/cfkit/env"
)

var (
	DefaultEurekaServiceName = "eureka"
	DefaultEurekaProtocol    = "http"
	DefaultEurekaPath        = "eureka"

	DefaultEurekaPort         = 80
	DefaultEurekaPollInterval = 30 * time.Second
	DefaultEurekaTimeout      = 10 * time.Second

	DefaultEurekaPortPropertyKey         = "port"
	DefaultEurekaTimeoutPropertyKey      = "timeout"
	DefaultEurekaPollIntervalPropertyKey = "poll_interval"
)

func Eureka() (*eureka.Client, error) {
	return EurekaWithName(DefaultEurekaServiceName)
}

func EurekaWithName(name string) (*eureka.Client, error) {
	svc, err := env.ServiceWithName(name)
	if err != nil {
		return nil, err
	}
	return eurekaLift(svc)
}

func EurekaWithTag(tag string) (*eureka.Client, error) {
	svc, err := env.ServiceWithTag(tag)
	if err != nil {
		return nil, err
	}
	return eurekaLift(svc)
}

var eurekaLift = EurekaFromService

func EurekaFromService(svc env.Service) (*eureka.Client, error) {
	uris, err := serviceURIs(svc)
	if err != nil {
		return nil, err
	}

	port := eurekaPort(svc)
	timeout := eurekaTimeout(svc)
	pollInterval := eurekaPollInterval(svc)

	return eureka.NewClient(uris, port, timeout, pollInterval), nil
}

func serviceURIs(svc env.Service) ([]string, error) {
	rawURIs, ok := svc.Credentials["uris"].([]interface{})
	if ok && len(rawURIs) > 0 {
		uris := make([]string, len(rawURIs))
		for i, raw := range rawURIs {
			uri, ok := raw.(string)
			if !ok {
				return []string{}, errors.New("Missing or invalid service URIs")
			}

			var err error
			if uri, err = augmentURI(uri); err != nil {
				return []string{}, err
			}
			uris[i] = uri
		}
		return uris, nil
	}

	uri, ok := svc.Credentials["uri"].(string)
	if !ok {
		return []string{}, errors.New("Missing or invalid service URIs")
	}

	var err error
	if uri, err = augmentURI(uri); err != nil {
		return []string{}, err
	}

	return []string{uri}, nil
}

func augmentURI(uri string) (string, error) {
	if uri == "" {
		return "", errors.New("Empty service URI")
	}

	if !strings.Contains(uri, "://") {
		uri = fmt.Sprintf("%s://%s", DefaultEurekaProtocol, uri)
	}

	url, err := url.Parse(uri)
	if err != nil {
		return "", fmt.Errorf("Error parsing service URI: %s", err)
	}

	if url.Path == "" || url.Path == "/" {
		url.Path = DefaultEurekaPath
	}

	return url.String(), nil
}

func eurekaPort(svc env.Service) int {
	if port, ok := svc.Credentials[DefaultEurekaPortPropertyKey].(int); ok {
		return port
	}
	return DefaultEurekaPort
}

func eurekaTimeout(svc env.Service) time.Duration {
	if timeout, ok := svc.Credentials[DefaultEurekaTimeoutPropertyKey].(int); ok {
		return time.Duration(timeout) * time.Second
	}
	return DefaultEurekaTimeout
}

func eurekaPollInterval(svc env.Service) time.Duration {
	if interval, ok := svc.Credentials[DefaultEurekaPollIntervalPropertyKey].(int); ok {
		return time.Duration(interval) * time.Second
	}
	return DefaultEurekaPollInterval
}
