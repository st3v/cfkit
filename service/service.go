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

func WithTag(tag string) (Service, error) {
	m, err := loadMap()
	if err != nil {
		return Service{}, err
	}

	return m.withTag(tag)
}

func WithName(name string) (Service, error) {
	m, err := loadMap()
	if err != nil {
		return Service{}, err
	}

	return m.withName(name)
}

type serviceMap map[string][]Service

func loadMap() (serviceMap, error) {
	jsonStr := os.Getenv(envVarName)
	if jsonStr == "" {
		return serviceMap{}, fmt.Errorf("%s not set", envVarName)
	}

	m := new(serviceMap)
	if err := json.Unmarshal([]byte(jsonStr), m); err != nil {
		return serviceMap{}, fmt.Errorf("Error parsing %s: %s", envVarName, err)
	}

	return *m, nil
}

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

func (m serviceMap) withName(name string) (Service, error) {
	for _, services := range m {
		for _, service := range services {
			if strings.ToUpper(name) == strings.ToUpper(service.Name) {
				return service, nil
			}
		}
	}
	return Service{}, fmt.Errorf("Service with name '%s' not found", name)
}
