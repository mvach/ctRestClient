package app

import (
	"bytes"
	"ctRestClient/logger"
	"fmt"
	"os/exec"
	"strings"
)

//counterfeiter:generate . KeepassCli
type KeepassCli interface {
	GetPassword(passwordName string) (string, error)
}

type keepassCli struct {
	dbFilePath string
	password   string
	logger     logger.Logger
}

func NewKeepassCli(dbFilePath string, password string, log logger.Logger) KeepassCli {
	return keepassCli{
		dbFilePath: dbFilePath,
		password:   password,
		logger:     log,
	}
}

func (s keepassCli) GetPassword(passwordName string) (string, error) {
	// Use the same command format that works on your command line
	cmd := exec.Command("keepassxc-cli", "show", "-q", "-a", "Password", s.dbFilePath, passwordName)
	cmd.Stdin = bytes.NewBufferString(s.password + "\n")

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command error: %v, stderr: %s, stdout: %s", err, stderr.String(), out.String())
	}

	// Trim whitespace from the output
	result := strings.TrimSpace(out.String())

	return result, nil
}
