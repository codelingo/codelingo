package service

import (
	"bufio"
	"errors"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/lingo-reviews/dev/tenet/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

type local struct {
	program    string
	args       []string
	socketAddr string
	process    *os.Process
}

func init() {
	grpclog.SetLogger(log.GetLogger())
}

// NewLocal allows you to run a program on the localhost as a micro-service.
func NewLocal(program string, args ...string) Service {
	log.Println("NewLocal service")
	return &local{
		program: program,
		args:    args,
	}
}

// StartService starts up the program as a micro-service server.
func (l *local) Start() error {

	// set a fixed socket and manually start process to help with debugging.
	// l.socketAddr = "@00208"
	// return nil

	// Start up the mirco-service
	log.Println("starting process")
	cmd := exec.Command(l.program, l.args...)
	p, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}

	// Get the socket address from the server.
	b := bufio.NewReader(p)
	line, _, err := b.ReadLine()
	if err != nil {
		log.Println("unable to get socket address from server, stopping tenet")
		l.Stop()
		return err
	}

	l.socketAddr = strings.TrimSuffix(string(line), "\n")
	l.process = cmd.Process

	return nil
}

// StopService stops the backing tenet server. If Common as a non-nil
// connection to the server, that will be closed.
func (l *local) Stop() (err error) {
	if l.process != nil {
		log.Println("killing process")
		if err = l.process.Kill(); err != nil {
			log.Fatalf("did not stop %s: %v", l.program, err)
		}
	}
	return
}

func (l *local) IsRunning() bool {
	return l.process.Signal(syscall.Signal(0)) == nil
}

func (l *local) DialGRPC() (*grpc.ClientConn, error) {
	if l.socketAddr == "" {
		return nil, errors.New("socket address is empty. Is the service started?")
	}
	log.Println("dialing server")
	return grpc.Dial(l.socketAddr, grpc.WithDialer(dialer()), grpc.WithInsecure())
}

func dialer() dialerFunc {
	switch runtime.GOOS {
	case "windows":
		return localWindowsDialer
	case "linux", "darwin", "freebsd":
		return localUnixDialer
	}
	return localUnixDialer
}

type dialerFunc func(addr string, timeout time.Duration) (net.Conn, error)

func localUnixDialer(addr string, timeout time.Duration) (net.Conn, error) {
	typ := "unix"
	raddr := net.UnixAddr{addr, typ}
	return net.DialUnix(typ, nil, &raddr)
}

func localWindowsDialer(addr string, timeout time.Duration) (net.Conn, error) {
	// TODO(waigani) implement
	panic("not implemented")
	return nil, nil
}
