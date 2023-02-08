package main

import (
	"fmt"
	"strings"

	"github.com/melbahja/goph"
)

func runCommand(client *goph.Client, cmd string) error {
	out, err := client.Run(cmd)
	if err != nil {
		trimOut := strings.TrimSuffix(string(out), "\n")
		return fmt.Errorf(
			"failed to run '%s', output: %s ,error: %v", cmd, trimOut, err,
		)
	}
	return nil
}

func runCommandOut(client *goph.Client, cmd string) (string, error) {
	out, err := client.Run(cmd)
	if err != nil {
		trimOut := strings.TrimSuffix(string(out), "\n")
		return "", fmt.Errorf(
			"failed to run '%s', output: %s ,error: %v", cmd, trimOut, err,
		)
	}
	return string(out), nil
}
