package binary

import (
	"bufio"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

type service struct {
	program    string
	args       []string
	socketAddr string
	process    *os.Process
}

func init() {
	grpclog.SetLogger(log.GetLogger())
}

// NewService allows you to run a program on the localhost as a micro-service.
func NewService(program string, args ...string) *service {
	log.Println("NewLocal service")
	return &service{
		program: program,
		args:    args,
	}
}

// StartService starts up the program as a micro-service server.
func (l *service) Start() error {

	// set a fixed socket and manually start process to help with debugging.
	// l.socketAddr = "@01c67"
	// return nil

	// Start up the mirco-service
	log.Println("starting process", l.program)
	log.Println(l.args)
	cmd := exec.Command(l.program, l.args...)
	p, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Trace(err)
	}
	if err := cmd.Start(); err != nil {
		return errors.Trace(err)
	}
	l.process = cmd.Process

	// Get the socket address from the server.
	b := bufio.NewReader(p)
	line, _, err := b.ReadLine()
	if err != nil {
		log.Println("unable to get socket address from server, stopping tenet")
		l.Stop()
		return errors.Trace(err)
	}

	l.socketAddr = strings.TrimSuffix(string(line), "\n")

	return nil
}

// StopService stops the backing tenet server. If Common as a non-nil
// connection to the server, that will be closed.
func (l *service) Stop() (err error) {
	if l.process != nil {
		log.Println("killing process")
		if err = l.process.Kill(); err != nil {
			log.Fatalf("did not stop %s: %v", l.program, err)
		}
	}
	return
}

// func (l *service) IsRunning() bool {
// 	return l.process.Signal(syscall.Signal(0)) == nil
// }

func (l *service) DialGRPC() (*grpc.ClientConn, error) {
	if l.socketAddr == "" {
		return nil, errors.New("socket address is empty. Is the service started?")
	}
	log.Println("dialing server")
	return grpc.Dial(l.socketAddr, grpc.WithDialer(dialer()), grpc.WithInsecure())
}

func dialer() dialerFunc {
	switch runtime.GOOS {
	case "windows":
		return serviceWindowsDialer
	case "linux", "freebsd":
		return serviceUnixDialer
	case "darwin":
		return serviceTcpDialer
	}
	return serviceUnixDialer
}

type dialerFunc func(addr string, timeout time.Duration) (net.Conn, error)

func serviceTcpDialer(addr string, timeout time.Duration) (net.Conn, error) {
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	return net.DialTCP("tcp", nil, raddr)
}

func serviceUnixDialer(addr string, timeout time.Duration) (net.Conn, error) {
	typ := "unix"
	raddr := net.UnixAddr{addr, typ}
	return net.DialUnix(typ, nil, &raddr)
}

func serviceWindowsDialer(addr string, timeout time.Duration) (net.Conn, error) {
	// TODO(waigani) implement
	panic("not implemented")
	return nil, nil
}
