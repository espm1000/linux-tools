package pkg

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
)

const dockerKeyRing string = "https://download.docker.com/linux/debian/gpg"
const dockerKey string = "/etc/apt/keyrings/docker.asc"
const dockerRepo string = "/etc/apt/sources.list.d/docker.sources"

type Dependencies struct {
	BasicDependencies []string `json:"basic_dependencies"`
	Docker            []string `json:"docker"`
	DebianDev         []string `json:"debian_dev"`
}

func loadDependencies() Dependencies {
	return GetPackageList()
}

func InstallAptDependencies(c *Config) error {
	slog.Info("installing updates")
	cmd := exec.Command("sudo", c.packageManager, "update", "-y")
	cmd.Stdin = os.Stdin
	if c.verbose {
		cmd.Stdout = os.Stdout
	}
	if err := cmd.Start(); err != nil {
		slog.Error("error running command", "error", err)
		return err
	}
	if err := cmd.Wait(); err != nil {
		slog.Error("error running command", "error", err)
		return err
	}
	slog.Info("complete.")
	return nil
}

func InstallDevTools(c *Config, verbose bool) error {
	if c.distro != "debian" {
		slog.Error("invalid operating system", "want", "debian", "have", c.distro)
		return errors.New("invalid OS")
	}
	var cmdList []exec.Cmd
	deps := loadDependencies()
	for _, dep := range deps.DebianDev {
		cmd := exec.Command("sudo", c.packageManager, "install", "-y", dep)
		cmdList = append(cmdList, *cmd)
	}
	for _, c := range cmdList {
		slog.Info("installing dependency", "package", c.Args[4])
		if verbose {
			c.Stderr = os.Stderr
		}
		if err := c.Run(); err != nil {
			return err
		}
	}
	return nil
}

func InstallDocker(c *Config) error {
	// Determine OS
	switch c.distro {
	case "debian":
		if err := installDockerDebian(c); err != nil {
			return err
		}
		return nil
	case "redhat":
		//installDockerRedhat(c)
		slog.Info("redhat currently noop")
		return nil
	default:
		fmt.Println("unknown OS")
	}
	return nil
}

func installDockerDebian(c *Config) error {
	var cmdList []*exec.Cmd
	var dockerCmdList []*exec.Cmd
	if c.currentUser != "root" {
		slog.Error("user must be root to install docker")
		return errors.New("non-root user")
	}
	slog.Info("reading packages list")
	deps := loadDependencies()

	for _, dep := range deps.BasicDependencies {
		slog.Info("preparing to install dependency", "package", dep)
		cmd := exec.Command("sudo", c.packageManager, "install", "-y", dep)
		cmdList = append(cmdList, cmd)
	}

	for _, cmd := range cmdList {
		slog.Info("installing dependency", "package", cmd.Args[4])
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	slog.Info("done.")
	slog.Info("setting docker apt key")
	if err := installDockerKeyringDebian(c); err != nil {
		slog.Error("error", "error", err)
		return err
	}
	slog.Info("adding docker apt repository")
	if err := writeDockerAptSource(c); err != nil {
		return err
	}
	if err := InstallAptDependencies(c); err != nil {
		return err
	}
	for _, dep := range deps.Docker {
		slog.Info("preparing to install dependency", "package", dep)
		cmd := exec.Command("sudo", c.packageManager, "install", "-y", dep)
		dockerCmdList = append(dockerCmdList, cmd)
	}
	for _, cmd := range dockerCmdList {
		slog.Info("installing dependency", "package", cmd.Args[4])
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}

func installDockerKeyringDebian(c *Config) error {
	if c.distro != "debian" {
		return fmt.Errorf("unsupported OS: %s", c.distro)
	}
	outFile, err := os.OpenFile(dockerKey, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := outFile.Close(); err != nil {
			slog.Error("error closing file stream", "error", err)
		}
	}()
	resp, err := http.Get(dockerKeyRing)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("error closing http stream", "error", err)
		}
	}()

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func writeDockerAptSource(c *Config) error {
	content := `Types: deb
URIs: https://download.docker.com/linux/debian
Suites: %s
Components: stable
Architectures: %s
Signed-By: %s`
	updated := fmt.Sprintf(content, c.OSInfo.VersionInfo, c.OSInfo.Arch, dockerKey)
	err := os.WriteFile(dockerRepo, []byte(updated), 0644)
	if err != nil {
		return err
	}

	return nil
}
