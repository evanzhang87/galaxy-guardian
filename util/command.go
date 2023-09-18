package util

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

// Command executes the named program with the given arguments. If it does not
// exit within timeout,it will be killed.
func Command(timeout time.Duration, name string, arg ...string) ([]byte, error) {
	// avoid leaking param
	var commandName = name[:]
	var args []string
	for _, a := range arg {
		args = append(args, a)
	}
	cmd := exec.Command(commandName, args...)
	randomBytes := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	cmd.Stdout = randomBytes
	cmd.Stderr = stderr
	done := make(chan string, 1)
	cmd.Start()
	timer := time.After(timeout)
	timeoutFlag := false
	go func(cmd *exec.Cmd) {
		for {
			select {
			case <-done:
				return
			case <-timer:
				timeoutFlag = true
				err := cmd.Process.Signal(os.Kill)
				log.Printf("ERROR: %v\n", err)
				return
			}
		}

	}(cmd)
	cmd.Wait()
	if timeoutFlag {
		return randomBytes.Bytes(), fmt.Errorf("%v", "timeout")
	}
	done <- "finish"
	errMsg := ""
	if stderr.Len() != 0 {
		errMsg += "," + stderr.String()
	}
	if len(errMsg) > 0 {
		return randomBytes.Bytes(), fmt.Errorf("%v", errMsg)
	} else {
		return randomBytes.Bytes(), nil
	}
}

// ReadCommand runs command name with args and calls line for each line from its
// stdout. Command is interrupted (if supported by Go) after 10 seconds and
// killed after 20 seconds.
func ReadCommand(line func(string) error, name string, arg ...string) error {
	return ReadCommandTimeout(time.Second*10, line, name, arg...)
}

// ReadCommandTimeout is the same as ReadCommand with a specifiable timeout.
func ReadCommandTimeout(timeout time.Duration, line func(string) error, name string, arg ...string) error {
	b, err := Command(timeout, name[:], arg...)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewBuffer(b))
	for scanner.Scan() {
		if err := line(scanner.Text()); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("ERROR: %v: %v\n", name[:], err)
	}
	return nil
}
