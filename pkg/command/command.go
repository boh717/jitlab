package command

import "os/exec"

type CommandClient interface {
	Run(command string, args ...string) ([]byte, error)
}

type CommandClientImpl struct{}

func (c CommandClientImpl) Run(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).CombinedOutput()
}
