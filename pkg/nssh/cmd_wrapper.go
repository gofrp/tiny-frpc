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
	name    string
	closeCh chan struct{}

	command string
	cmd     *exec.Cmd
}

func NewCmdWrapper(ctx context.Context, command string, closeCh chan struct{}) *CmdWrapper {
	parts := strings.Fields(command)

	wrapper := &CmdWrapper{
		cmd:     exec.CommandContext(ctx, parts[0], parts[1:]...),
		command: command,
		closeCh: closeCh,
	}

	go wrapper.wait()

	return wrapper
}

func (cs *CmdWrapper) wait() {
	<-cs.closeCh
	cs.cmd.Wait()
}

func (cs *CmdWrapper) ExecuteCommand(ctx context.Context) {
	outputCh := make(chan string)
	errCh := make(chan error, 1)
	defer close(outputCh)
	defer close(errCh)

	go func() {
		for out := range outputCh {
			// do not use log, use standard print to better show output
			fmt.Println(out)
		}
	}()

	go func() {
		for err := range errCh {
			log.Errorf("run cmd: %v error: %v", cs.command, err)
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
		errCh <- err
		return
	}

	stdoutReader := bufio.NewReader(stdoutPipe)
	stderrReader := bufio.NewReader(stderrPipe)

	go cs.readPipe(stdoutReader, outputCh)
	go cs.readPipe(stderrReader, outputCh)

	if err := cs.cmd.Wait(); err != nil {
		errCh <- err
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
