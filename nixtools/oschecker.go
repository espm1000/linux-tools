package main

import (
	"bufio"
	"errors"
	"log/slog"
	"os"
	"runtime"
	"strings"
)

func checkOS() (*Config, error) {
	cfg := &Config{}
	os := runtime.GOOS
	switch os {
	case "linux":
		slog.Info("operating system detected", "os", os)
		config, err := getLinuxDistro()
		if err != nil {
			return nil, err
		}
		cfg = config
	default:
		slog.Info("failed to detect OS")
	}
	return &Config{
		os:             os,
		distro:         cfg.distro,
		packageManager: cfg.packageManager,
	}, nil
}

func getLinuxDistro() (*Config, error) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("error closing io stream", "error", err)
		}
	}()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		_, found := strings.CutPrefix(line, "PRETTY_NAME=")

		if found {
			switch {
			case strings.Contains(normalizeString(line), "debian") || strings.Contains(normalizeString(line), "ubuntu"):
				return &Config{
					packageManager: "apt",
					distro:         "debian",
				}, nil
			case strings.Contains(normalizeString(line), "redhat") || strings.Contains(normalizeString(line), "fedora"):
				return &Config{
					packageManager: "dnf",
					distro:         "redhat",
				}, nil
			default:
				slog.Error("unable to detect distribution", "response", line)
				return nil, errors.New("no distro detected")
			}
		}
	}
	return nil, nil
}

func normalizeString(s string) string {
	return strings.ToLower(s)
}
