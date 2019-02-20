package gdb

import (
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/solo-io/squash/pkg/debuggers"

	"fmt"
	"time"
)

type GdbInterface struct{}

type gdbDebugServer struct {
	port int
	cmd  *exec.Cmd
}

func (g *gdbDebugServer) Detach() error {
	g.cmd.Process.Signal(syscall.SIGINT)
	return nil
}

func (g *gdbDebugServer) Port() int {
	return g.port
}

func (g *gdbDebugServer) HostType() debuggers.DebugHostType {
	return debuggers.DebugHostTypeClient
}

func (d *gdbDebugServer) Cmd() *exec.Cmd {
	return nil
}

func (g *GdbInterface) Attach(pid int) (debuggers.DebugServer, error) {

	log.WithField("pid", pid).Debug("AttachToLiveSession called")
	cmd := exec.Command("gdbserver", "--attach", ":0", fmt.Sprintf("%d", pid))
	cmd.Start()
	log.Debug("starting gdbserver for user started, trying to get port")
	time.Sleep(time.Second)
	port, err := debuggers.GetPort(cmd.Process.Pid)
	if err != nil {
		log.WithField("err", err).Error("can't get gdbserver port")
		cmd.Process.Kill()
		cmd.Process.Release()
		return nil, err
	}

	// be polite and wait
	go cmd.Wait()
	gds := &gdbDebugServer{
		port: port,
		cmd:  cmd,
	}
	return gds, nil
}
