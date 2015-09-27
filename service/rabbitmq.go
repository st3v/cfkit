package service

import (
	"errors"
	"strings"

	"github.com/st3v/cfkit/env"
	"github.com/streadway/amqp"
)

const DefaultRabbitTag = "rabbitmq"

var amqpDialer = amqp.Dial

type RabbitMQ struct {
	uri string
}

func (r *RabbitMQ) Dial() (*amqp.Connection, error) {
	return amqpDialer(r.uri)
}

func (r *RabbitMQ) URI() string {
	return r.uri
}

func Rabbit() (*RabbitMQ, error) {
	return RabbitWithTag(DefaultRabbitTag)
}

func RabbitWithTag(tag string) (*RabbitMQ, error) {
	return find(env.ServiceWithTag, tag)
}

func RabbitWithName(name string) (*RabbitMQ, error) {
	return find(env.ServiceWithName, name)
}

func RabbitFromService(svc env.Service) (*RabbitMQ, error) {
	uri, ok := svc.Credentials["uri"].(string)
	if !ok || !strings.HasPrefix(uri, "amqp://") {
		return nil, errors.New("Invalid AMQP URI")
	}
	return &RabbitMQ{uri}, nil
}

var serviceLift = RabbitFromService

type lookupFn func(string) (env.Service, error)

func find(lookup lookupFn, id string) (*RabbitMQ, error) {
	svc, err := lookup(id)
	if err != nil {
		return nil, err
	}
	return serviceLift(svc)
}
