package cfhelper

import (
	"log"
	"os"
)

func Domain() string {
	return getenv("CF_DOMAIN")
}

func API() string {
	return getenv("CF_API")
}

func Username() string {
	return getenv("CF_USERNAME")
}

func Password() string {
	return getenv("CF_PASSWORD")
}

func Org() string {
	return getenv("CF_ORG")
}

func RabbitServiceName() string {
	return getenv("CF_RABBIT_SERVICE_NAME")
}

func RabbitServicePlan() string {
	return getenv("CF_RABBIT_SERVICE_PLAN")
}

func getenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Env var %s not defined", key)
	}
	return val
}
