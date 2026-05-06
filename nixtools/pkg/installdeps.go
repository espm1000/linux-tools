package pkg

import (
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
