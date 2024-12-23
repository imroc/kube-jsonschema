package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"time"
)

var regAddr = regexp.MustCompile(`127.0.0.1:(\d+)`)

func runKubectlProxy() (cmd *exec.Cmd, port string, err error) {
	cmd = exec.Command("kubectl", "proxy", "--port=0")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return
	}
	stdoutReader := bufio.NewReader(stdout)
	exit := make(chan error)
	timeout := time.After(2 * time.Second)
	if err = cmd.Start(); err != nil {
		return
	}
	go func() {
		exit <- cmd.Wait()
	}()

	output := make(chan string)

	go func() {
		line, _, _ := stdoutReader.ReadLine()
		output <- string(line)
	}()

	select {
	case <-timeout:
		err = errors.New("kubectl proxy executed timeout")
		return
	case e := <-exit:
		if e != nil {
			err = e
		} else {
			stderrBytes, e := io.ReadAll(stderr)
			if e != nil {
				err = e
				return
			}
			if len(stderrBytes) > 0 {
				err = errors.New(string(stderrBytes))
				return
			}
			err = errors.New("kubectl proxy exited unexpectedly")
		}
		return
	case line := <-output:
		if len(line) == 0 {
			err = errors.New("kubectl proxy output nothing")
			return
		}
		ss := regAddr.FindStringSubmatch(line)
		if len(ss) != 2 {
			err = fmt.Errorf("kubectl proxy output unknown format: %s", line)
			return
		}
		port = ss[1]
		return
	}
}
