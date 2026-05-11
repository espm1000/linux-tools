package pkg

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

const osRelease string = "/etc/os-release"

type OSInfo struct {
	VersionInfo string
	Arch        string
}

type Config struct {
	OSInfo         *OSInfo
	currentUser    string
	hostname       string
	homeDiretory   string
	packageManager string
	distro         string
	os             string
	verbose        bool
}

func GenerateConfig() (*Config, error) {
	cu, err := getCurrentUser()
	if err != nil {
		slog.Error("error getting user", "error", err)
		return nil, err
	}

	hostname, err := getHostname()
	if err != nil {
		return nil, err
	}

	var homeDir string
	if cu != "root" {
		homeDir, err = checkEnvironmentFile(cu)
		if err != nil {
			return nil, err
		}
	}
	cfg, err := checkOS()
	if err != nil {
		return nil, err
	}

	version, err := getOSDetails()
	if err != nil {
		return nil, err
	}

	arch, err := getArch(cfg)
	if err != nil {
		return nil, err
	}

	return &Config{
		currentUser:    cu,
		hostname:       hostname,
		homeDiretory:   homeDir,
		os:             cfg.os,
		distro:         cfg.distro,
		packageManager: cfg.packageManager,
		verbose:        false,
		OSInfo: &OSInfo{
			VersionInfo: version,
			Arch:        strings.Replace(arch, "\n", "", 1),
		},
	}, nil

}

func getCurrentUser() (string, error) {
	slog.Info("checking for user settings")
	current_user, err := user.Current()
	if err != nil {
		slog.Error("error reading local user", "error", err)
		return "", err
	}
	slog.Info("detected username", "username", current_user.Username)
	return current_user.Username, nil
}

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
	file, err := os.Open(osRelease)
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

func parseOSReleaseDetails(key string) (string, error) {
	data, err := os.ReadFile(osRelease)
	if err != nil {
		return "", err
	}
	prefix := key + "="

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			value := strings.TrimPrefix(line, prefix)
			return strings.Trim(value, `"'`), nil
		}
	}
	return "", fmt.Errorf("%s not found", prefix)
}

func getOSDetails() (string, error) {
	version, err := parseOSReleaseDetails("VERSION_CODENAME")
	if err != nil {
		return "", err
	}
	return version, nil
}

func getArch(c *Config) (string, error) {
	if c.distro == "debian" {
		arch, err := exec.Command("dpkg", "--print-architecture").Output()
		if err != nil {
			return "", err
		}
		return string(arch), nil
	}
	return "unknown", nil
}

func normalizeString(s string) string {
	return strings.ToLower(s)
}
