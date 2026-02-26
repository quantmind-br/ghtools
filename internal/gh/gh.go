package gh

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

func Run(args ...string) (string, error) {
	cmd := exec.Command("gh", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%s: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func RunJSON(result interface{}, args ...string) error {
	out, err := Run(args...)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(out), result)
}

func CheckAuth() error {
	cmd := exec.Command("gh", "auth", "status")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not authenticated with GitHub CLI. Run: gh auth login")
	}
	return nil
}

func CheckInstalled() error {
	_, err := exec.LookPath("gh")
	if err != nil {
		return fmt.Errorf("gh CLI not found. Install from https://cli.github.com")
	}
	return nil
}
