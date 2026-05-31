package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
)

type Config struct {
	Namespace            string
	DeploymentName       string
	ServiceName          string
	DefaultAdminPassword string
	PodName              string
	Kpath                string
}

func (c Config) GetSecret() (string, error) {
	out, err := exec.Command("kubectl", "get", "secret", "--namespace", c.Namespace, "grafana", "-o", "jsonpath=\"{.data.admin-password}\"").Output()
	if err != nil {
		slog.Error("error running command", "error", err)
		return "", err
	}
	encodedValue, err := strconv.Unquote(string(out))
	if err != nil {
		slog.Error("error removing quotes", "error", err)
		return "", err
	}
	decodedValue, err := base64.StdEncoding.DecodeString(encodedValue)
	if err != nil {
		slog.Error("error decoding string", "error", err)
		return "", err
	}
	return string(decodedValue), nil
}

func (c Config) GetPodName() (string, error) {
	podName, err := exec.Command("kubectl", "get", "pods", "--namespace", c.Namespace, "-l", c.Kpath, "-o", "jsonpath=\"{.items[0].metadata.name}\"").Output()
	if err != nil {
		slog.Error("error getting pod name", "error", err)
		return "", err
	}
	return string(podName), nil
}

func (c Config) GetConfig() (*Config, error) {
	pass, err := c.GetSecret()
	if err != nil {
		return nil, err
	}
	podName, err := c.GetPodName()
	if err != nil {
		return nil, err
	}
	return &Config{
		DefaultAdminPassword: string(pass),
		PodName:              podName,
	}, nil
}

func main() {
	k := Config{
		Namespace: "default",
		Kpath:     "app.kubernetes.io/name=grafana,app.kubernetes.io/instance=grafana",
	}
	cfg, err := k.GetConfig()
	if err != nil {
		slog.Error("error occurred", "error", err)
		panic(err)
	}
	fmt.Printf("admin password: %s\n", cfg.DefaultAdminPassword)
	fmt.Printf("pod name: %s\n", cfg.PodName)
}
