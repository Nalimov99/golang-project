package databasetest

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"testing"
)

// container tracks information about docker container started for test
type container struct {
	Name  string
	Ports string //IP:Port
}

const containerName = "gotest"

// startContainer runs a postgres container to execute commands
func startContainer(t *testing.T) *container {
	t.Helper()

	cmd := exec.Command("docker", "run",
		"--name", containerName,
		"-d",
		"-e", "POSTGRES_PASSWORD=postgres",
		"-e", "POSTGRESQL_USER=postgres",
		"-p", "5433:5432",
		"postgres:14.1-alpine")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		t.Fatalf("could not start container: %s. Error: %v", containerName, err)
	}

	out.Reset()
	cmd = exec.Command("docker", "inspect", containerName)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect docker container %v", err)
	}

	var doc []struct {
		NetworkSettings struct {
			Ports struct {
				TCP5432 []struct {
					HostIp   string `json:"HostIp"`
					HostPort string `json:"HostPort"`
				} `json:"5432/tcp"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}

	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json %v", err)
	}

	network := doc[0].NetworkSettings.Ports.TCP5432[0]

	return &container{
		Name:  containerName,
		Ports: network.HostIp + ":" + network.HostPort,
	}
}

// stopContainer stop and removes the specified container
func stopContainer(t *testing.T) {
	t.Helper()

	if err := exec.Command("docker", "stop", containerName).Run(); err != nil {
		t.Fatalf("could not stop container: %s. Error: %s", containerName, err)
	}

	t.Logf("container stopped: %s", containerName)

	if err := exec.Command("docker", "rm", containerName).Run(); err != nil {
		t.Fatalf("could not remove container: %s. Error: %s", containerName, err)
	}

	t.Logf("container removed: %s", containerName)
}
