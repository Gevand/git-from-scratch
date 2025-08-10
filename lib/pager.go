package lib

import (
	"io"
	"os"
	"os/exec"
)

const PAGER_CMD = "less"

type Pager struct {
	StdIn io.WriteCloser
	Cmd   *exec.Cmd
}

func NewPager() *Pager {
	return &Pager{}
}

func (pager *Pager) Initialize() error {
	reader, writer := io.Pipe()
	pager.StdIn = writer
	cmd := exec.Command(PAGER_CMD, "-c")
	pager.Cmd = cmd
	cmd.Stdin = reader
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	cmd.Start()
	return nil
}
