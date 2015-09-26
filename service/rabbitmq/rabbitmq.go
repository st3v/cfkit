package rabbitmq

import (
	"github.com/st3v/cfkit/service"
	"github.com/streadway/amqp"
)

const DefaultTag = "rabbitmq"

var dialer = amqp.Dial

type RabbitService interface {
	Dial() (*amqp.Connection, error)
}

type rabbit struct {
	uri string
}

func (r *rabbit) Dial() (*amqp.Connection, error) {
	return dialer(r.uri)
}

func Service() (*rabbit, error) {
	return ServiceWithTag(DefaultTag)
}

func ServiceWithTag(tag string) (*rabbit, error) {
	srv, err := service.WithTag(tag)
	if err != nil {
		return nil, err
	}

	uri, err := URI(srv)
	if err != nil {
		return nil, err
	}

	return &rabbit{uri}, nil
}
