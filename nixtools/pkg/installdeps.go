package pkg

import (
	"log/slog"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
)

type UpdateConfig struct {
	packageManager string
}

func (c *Config) installAptDependencies(verbose string) error {
	slog.Info("installing updates")
	cmd := exec.Command("sudo", c.packageManager, "update", "-y")
	cmd.Stdin = os.Stdin
	if verbose == "true" {
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

func initSpinner() *spinner.Spinner {
	return spinner.New(spinner.CharSets[9], 100*time.Millisecond)
}
