package rabbitmq

import (
	"errors"
	"strings"

	"github.com/st3v/cfkit/service"
)

var URI = func(s service.Service) (string, error) {
	switch s.Label {
	case "p-rabbitmq":
		return pRabbitURI(s)
	default:
		return simpleURI(s)
	}
}

func pRabbitURI(s service.Service) (string, error) {
	protos, ok := s.Credentials["protocols"].(map[string]interface{})
	if !ok {
		return "", errors.New("Invalid service credentials")
	}

	amqpProto, ok := protos["amqp"].(map[string]interface{})
	if !ok {
		return "", errors.New("Invalid AMQP protocol credentials")
	}

	uris, ok := amqpProto["uris"].([]interface{})
	if !ok || len(uris) < 1 {
		return "", errors.New("Invalid AMQP URIs")
	}

	uri, ok := uris[0].(string)
	if !ok || uri == "" {
		return "", errors.New("Invalid AMQP URI")
	}

	return uri, nil
}

func simpleURI(s service.Service) (string, error) {
	uri, ok := s.Credentials["uri"].(string)
	if !ok || !strings.HasPrefix(uri, "amqp://") {
		return "", errors.New("Invalid AMQP URI")
	}
	return uri, nil
}
