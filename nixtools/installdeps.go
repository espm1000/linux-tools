package main

import (
	"log/slog"
	"os"
	"os/exec"
)

func install_apt_dependencies() error {
	slog.Info("checking for updates")
	cmd := exec.Command("sudo", "apt", "update", "-y")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		slog.Error("error running command", "error", err)
		return err
	}
	slog.Info("complete.")
	return nil
}
