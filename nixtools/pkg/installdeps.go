package pkg

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
)

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
