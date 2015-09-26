package rabbitmq

import (
	"errors"
	"strings"

	"github.com/st3v/cfkit/service"
	"github.com/streadway/amqp"
)

const DefaultTag = "rabbitmq"

var dialer = amqp.Dial

type Svc struct {
	uri string
}

func (s *Svc) Dial() (*amqp.Connection, error) {
	return dialer(s.uri)
}

func (s *Svc) URI() string {
	return s.uri
}

func Service() (*Svc, error) {
	return ServiceWithTag(DefaultTag)
}

func ServiceWithTag(tag string) (*Svc, error) {
	return find(service.WithTag, tag)
}

func ServiceWithName(name string) (*Svc, error) {
	return find(service.WithName, name)
}

type lookupFn func(string) (service.Service, error)

func find(lookup lookupFn, id string) (*Svc, error) {
	svc, err := lookup(id)
	if err != nil {
		return nil, err
	}
	return newFromService(svc)
}

var newFromService = func(svc service.Service) (*Svc, error) {
	switch svc.Label {
	case "p-rabbitmq":
		return prabbit(svc.Credentials)
	default:
		return standard(svc.Credentials)
	}
}

type credentials map[string]interface{}

func prabbit(creds credentials) (*Svc, error) {
	protos, ok := creds["protocols"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid service credentials")
	}

	amqpProto, ok := protos["amqp"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid AMQP protocol credentials")
	}

	return standard(amqpProto)
}

func standard(creds credentials) (*Svc, error) {
	uri, ok := creds["uri"].(string)
	if !ok || !strings.HasPrefix(uri, "amqp://") {
		return nil, errors.New("Invalid AMQP URI")
	}
	return &Svc{uri}, nil
}
