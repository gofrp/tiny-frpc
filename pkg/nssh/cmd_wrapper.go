package nssh

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/gofrp/tiny-frpc/pkg/util/log"
)

type CmdWrapper struct {
	command string
	cmd     *exec.Cmd

	outputCh chan string
	errCh    chan error
}

func NewCmdWrapper(ctx context.Context, command string) *CmdWrapper {
	parts := strings.Fields(command)

	wrapper := &CmdWrapper{
		cmd:     exec.CommandContext(ctx, parts[0], parts[1:]...),
		command: command,

		outputCh: make(chan string),
		errCh:    make(chan error, 1),
	}

	return wrapper
}

func (cs *CmdWrapper) ExecuteCommand(ctx context.Context) {
	go func() {
		for out := range cs.outputCh {
			// do not use log, use standard print to better show output
			fmt.Println(out)
		}
	}()

	go func() {
		for err := range cs.errCh {
			log.Errorf("run cmd: [%v] get error: %v", cs.command, err)
		}
	}()

	stdoutPipe, err := cs.cmd.StdoutPipe()
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		close(errCh)
		return
	}

	stderrPipe, err := cs.cmd.StderrPipe()
	if err != nil {
		errCh := make(chan error, 1)
		errCh <- err
		close(errCh)
		return
	}

	if err := cs.cmd.Start(); err != nil {
		cs.errCh <- err
		return
	}

	stdoutReader := bufio.NewReader(stdoutPipe)
	stderrReader := bufio.NewReader(stderrPipe)

	go cs.readPipe(stdoutReader, cs.outputCh)
	go cs.readPipe(stderrReader, cs.outputCh)

	if err := cs.cmd.Wait(); err != nil {
		cs.errCh <- err
	}
}

func (cs *CmdWrapper) readPipe(pipe *bufio.Reader, outputCh chan<- string) {
	for {
		line, err := pipe.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				outputCh <- fmt.Sprintf("Error: %s", err)
			}
			break
		}
		outputCh <- line
	}
}

func (cs *CmdWrapper) Close() {
	if cs.cmd != nil && cs.cmd.Process != nil {
		cs.cmd.Process.Kill()
	}
}
