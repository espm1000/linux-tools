package main

import (
	"log/slog"
	"os/exec"
)

func install_apt_dependencies() error {
	cmd := exec.Command("sudo", "apt", "update", "-y")
	if err := cmd.Run(); err != nil {
		slog.Error("error running command", "error", err)
		return err
	}
	return nil
}