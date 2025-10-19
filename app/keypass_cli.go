package app

import (
	"bytes"
	"ctRestClient/logger"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//counterfeiter:generate . KeepassCli
type KeepassCli interface {
	GetPassword(passwordName string) (string, error)
	IsPasswordValid(passwordName string) (bool, error)
}

type keepassCli struct {
	dbFilePath string
	password   string
	logger     logger.Logger
}

func NewKeepassCli(dbFilePath string, password string, log logger.Logger) (KeepassCli, error) {
	info, err := os.Stat(dbFilePath)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("the Keepass DB file '%s' could not be found", dbFilePath)
	}
	if err != nil {
		return nil, fmt.Errorf("error checking Keepass DB file '%s': %v", dbFilePath, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("the Keepass DB file '%s' exists but is not a regular file", dbFilePath)
	}
	return keepassCli{
		dbFilePath: dbFilePath,
		password:   password,
		logger:     log,
	}, nil
}

func (s keepassCli) GetPassword(passwordName string) (string, error) {
	cmd := exec.Command("keepassxc-cli", "show", "-q", "-a", "Password", s.dbFilePath, passwordName)
	cmd.Stdin = bytes.NewBufferString(s.password + "\n")

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("command error: %v, stderr: %s, stdout: %s", err, stderr.String(), out.String())
	}

	result := strings.TrimSpace(out.String())

	return result, nil
}

func (s keepassCli) IsPasswordValid(passwordName string) (bool, error) {
	cmd := exec.Command("keepassxc-cli", "ls", "-q", s.dbFilePath)
	cmd.Stdin = bytes.NewBufferString(s.password + "\n")

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err == nil {
		// exit code 0
		return true, nil
	}

	// If command executed but returned non-zero exit code, keepassxc-cli will
	// return an *exec.ExitError. Treat that as "password invalid" (false, nil).
	if _, ok := err.(*exec.ExitError); ok {
		return false, nil
	}

	return false, fmt.Errorf("command error: %v, stderr: %s, stdout: %s", err, stderr.String(), out.String())
}