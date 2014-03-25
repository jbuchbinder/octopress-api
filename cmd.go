package main

import (
	"bytes"
	"errors"
	"io"
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

	stdout, err := proc.StdoutPipe()
	if err != nil {
		return "", errors.New("Failed to create stdout pipe")
	}
	stderr, err := proc.StderrPipe()
	if err != nil {
		return "", errors.New("Failed to create stderr pipe")
	}

	err = proc.Start()
	if err != nil {
		return "", errors.New("Failed to start process : " + err.Error())
	}

	// Hack to write everything back to the user
	b := bytes.NewBufferString("")
	go io.Copy(b, stdout)
	go io.Copy(b, stderr)

	err = proc.Wait()
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
