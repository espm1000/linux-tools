package pkg

import (
	"bufio"
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
	deps := []string{"build-essential", "checkinstall", "libz-dev", "dh-make", "libssl-dev", "devscripts"}
	for _, dep := range deps {
		cmd := exec.Command("sudo", c.packageManager, "install", "-y", dep)
		cmdList = append(cmdList, *cmd)
	}
	for _, c := range cmdList {
		fmt.Printf("installing update %v\n", c.Args[4])
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
		return nil
	default:
		fmt.Println("unknown OS")
	}
	return nil
}

func installDockerDebian(c *Config) error {
	if c.currentUser != "root" {
		slog.Error("user must be root to install docker")
		return errors.New("non-root user")
	}
	slog.Info("reading packages list")
	packages, err := os.Open("./internal/templates/package.list")
	if err != nil {
		slog.Error("error reading file", "error", err)
		return err
	}
	defer func() error {
		if err := packages.Close(); err != nil {
			return err
		}
		return nil
	}()

	scanner := bufio.NewScanner(packages)
	var cmdList []*exec.Cmd
	for scanner.Scan() {
		slog.Info("installing package", "package", scanner.Text())
		cmd := exec.Command("sudo", c.packageManager, "install", "-y", scanner.Text())
		cmdList = append(cmdList, cmd)
	}
	for _, cmd := range cmdList {
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
	defer outFile.Close()
	resp, err := http.Get(dockerKeyRing)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func writeDockerAptSource(c *Config) error {
	tmpl, err := os.ReadFile("./internal/templates/docker.repo")
	if err != nil {
		slog.Error("error reading file", "error", err)
		return err
	}
	updated := fmt.Sprintf(string(tmpl), c.OSInfo.VersionInfo, c.OSInfo.Arch, dockerKey)
	err = os.WriteFile(string(tmpl), []byte(updated), 0644)
	if err != nil {
		return err
	}

	fmt.Println(updated)
	return nil
}
