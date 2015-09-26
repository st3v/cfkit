package cfhelper

import (
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"strings"
)

const cli = "cf"

func TargetOrg(org string) error {
	cmd := exec.Command(cli, "target", "-o", org)
	return runCommand(cmd, "Error targeting org")
}

func CreateSpace(space string) error {
	cmd := exec.Command(cli, "create-space", space)
	return runCommand(cmd, "Error creating space")
}

func TargetSpace(space string) error {
	cmd := exec.Command(cli, "target", "-s", space)
	return runCommand(cmd, "Error targeting space")
}

func DeleteSpace(space string) error {
	cmd := exec.Command(cli, "delete-space", "-f", space)
	return runCommand(cmd, "Error deleting space")
}

func Login(api, username, password string) error {
	cmd := exec.Command(cli, "api", api)
	if err := runCommand(cmd, "Error setting CF API endpoint"); err != nil {
		return err
	}

	cmd = exec.Command(cli, "login", "-u", username, "-p", password)
	return runCommand(cmd, "Error logging into CF")
}

func PushAppManifest(app, manifestPath string) error {
	cmd := exec.Command(cli, "push", app, "-f", manifestPath, "--no-start")
	return runCommand(cmd, "Error pushing app to CF")
}

func StartApp(app string) error {
	cmd := exec.Command(cli, "start", app)
	return runCommand(cmd, "Error starting app")
}

func DeleteApp(app string) error {
	cmd := exec.Command(cli, "delete", "-f", "-r", app)
	return runCommand(cmd, "Error deleting app")
}

func CreateService(serviceName, servicePlan, instanceID string) error {
	cmd := exec.Command(cli, "create-service", serviceName, servicePlan, instanceID)
	return runCommand(cmd, "Error creating service")
}

func BindService(appName, instanceID string) error {
	cmd := exec.Command(cli, "bind-service", appName, instanceID)
	return runCommand(cmd, "Error binding service")
}

func DeleteService(instanceID string) error {
	cmd := exec.Command(cli, "delete-service", "-f", instanceID)
	return runCommand(cmd, "Error deleting service")
}

func RandomServiceID() string {
	return fmt.Sprintf("service-instance-%s", randID())
}

func RandomSpaceName() string {
	return fmt.Sprintf("space-%s", randID())
}

func RandomAppName() string {
	return fmt.Sprintf("test-app-%s", randID())
}

func sanitize(str string) string {
	str = strings.Replace(str, Username(), "HIDDEN_USERNAME", -1)
	return strings.Replace(str, Password(), "HIDDEN_PASSWORD", -1)
}

func runCommand(cmd *exec.Cmd, errMsg string) error {
	log.Println(sanitize(strings.Join(cmd.Args, " ")))

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", errMsg, sanitize(string(out)))
	}

	log.Println(sanitize(string(out)))

	return nil
}

func randID() string {
	chars := []byte("abcdefghijklmnopqrstuvwxyz0123456789")
	result := make([]byte, 7)
	for i, _ := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
