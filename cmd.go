package main

import (
	"errors"
	"log"
	"os/exec"
)

// RunCmd is a simplistic method of running a local command with arguments.
// The cmd argument must be a fully qualified path to the destination
// executable file, otherwise it will not function.
func RunCmd(workingDirectory string, cmd string, argv []string) (string, error) {
	log.Print("Executing command : " + cmd + " in dir " + workingDirectory)

	proc := exec.Command(cmd, argv...)
	proc.Dir = workingDirectory

	b, err := proc.CombinedOutput()
	if err != nil {
		return "", errors.New("Failed to start process : " + err.Error())
	}

	return string(b), nil
}
