package service

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const envVarName = "VCAP_SERVICES"

type Service struct {
	Name        string                 `json:"name"`
	Label       string                 `json:"label"`
	Tags        []string               `json:"tags"`
	Plan        string                 `json:"plan"`
	Credentials map[string]interface{} `json:"credentials"`
}

type serviceMap map[string][]Service

func (m serviceMap) withTag(tag string) (Service, error) {
	for _, services := range m {
		for _, service := range services {
			for _, t := range service.Tags {
				if strings.ToUpper(t) == strings.ToUpper(tag) {
					return service, nil
				}
			}
		}
	}
	return Service{}, fmt.Errorf("Service with tag '%s' not found", tag)
}

func WithTag(tag string) (Service, error) {
	jsonStr := os.Getenv(envVarName)
	if jsonStr == "" {
		return Service{}, fmt.Errorf("%s not set", envVarName)
	}

	m := new(serviceMap)
	if err := json.Unmarshal([]byte(jsonStr), m); err != nil {
		return Service{}, err
	}

	return m.withTag(tag)
}
