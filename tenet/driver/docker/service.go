package docker

import (
	"net"
	"os"
	"os/exec"
	"time"

	goDocker "github.com/fsouza/go-dockerclient"
	"github.com/juju/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	"github.com/lingo-reviews/dev/tenet/log"
	"github.com/lingo-reviews/lingo/util"
)

type service struct {
	ip           string
	port         string
	process      *os.Process
	image        string
	containerID  string
	dockerClient *goDocker.Client
}

func init() {
	grpclog.SetLogger(log.GetLogger())
}

// TODO(waigani) don't pass client in - make it
func NewService(tenetName string) (*service, error) {
	return &service{
		image: tenetName,
	}, nil
}

func (s *service) Start() error {
	log.Print("docker.service.Start")
	c, err := s.client()
	if err != nil {
		return errors.Trace(err)
	}
	// dockerArgs := []string{"start", containerName}

	// TODO(waigani) check that pwd is correct when a tenet is started for a
	// subdir.
	pwd, err := os.Getwd()
	if err != nil {
		return errors.Trace(err)
	}

	// Start up the mirco-service

	internalPort := "8000/tcp"
	dockerPort := goDocker.Port(internalPort)

	// start a new container
	opts := goDocker.CreateContainerOptions{
		Config: &goDocker.Config{
			Image: s.image,
			ExposedPorts: map[goDocker.Port]struct{}{
				dockerPort: {}},
			AttachStdin: true,
			Tty:         true,
		},
		HostConfig: &goDocker.HostConfig{
			PublishAllPorts: true,
			Binds:           []string{pwd + ":/source:ro"},
			PortBindings: map[goDocker.Port][]goDocker.PortBinding{
				dockerPort: []goDocker.PortBinding{{
					HostIP:   "127.0.0.1",
					HostPort: "0",
				}},
			},
		},
	}

	container, err := c.CreateContainer(opts)
	if err != nil {
		return errors.Trace(err)
	}
	s.containerID = container.ID

	err = c.StartContainer(container.ID, nil)
	if err != nil {
		return errors.Annotatef(err, "error starting container %s", container.Name)
	}

	for container.NetworkSettings == nil {
		time.Sleep(1 * time.Microsecond)
		container, err = c.InspectContainer(container.ID)
		if err != nil {
			return errors.Trace(err)
		}
	}

	log.Printf("waiting for ports to bind for docker container", container.ID)
	var breakLoop bool
	go func() {
		<-time.After(5 * time.Second)
		breakLoop = true
	}()
	for container.NetworkSettings.Ports[dockerPort] == nil && !breakLoop {
		time.Sleep(1 * time.Microsecond)
		container, err = c.InspectContainer(container.ID)
		if err != nil {
			return errors.Trace(err)
		}
	}
	if breakLoop {
		return errors.New("timed out waiting for docker ports to bind")
	}

	log.Printf("%#v", container.NetworkSettings.Ports[dockerPort])

	ports := container.NetworkSettings.Ports[dockerPort]
	s.ip = ports[0].HostIP
	s.port = ports[0].HostPort

	log.Print("got to end of docker service.Start, no errors")
	return nil
}

func (s *service) client() (*goDocker.Client, error) {
	if s.dockerClient == nil {
		var err error
		s.dockerClient, err = util.DockerClient()
		if err != nil {
			return nil, errors.Trace(err)
		}
	}

	return s.dockerClient, nil
}

func (s *service) stop() {
	// Use exec so it's quick
	cmd := exec.Command("docker", "rm", "-f", s.containerID)
	if err := cmd.Start(); err != nil {
		log.Println("ERROR stopping tenet:", err)
		time.Sleep(1 * time.Microsecond)
		log.Println("trying to stop tenet again")
		s.stop()
		return
	}
	s.release(cmd)
}

func (s *service) release(cmd *exec.Cmd) {
	if err := cmd.Process.Release(); err != nil {
		log.Println("ERROR releasing process:", err)
		time.Sleep(1 * time.Microsecond)
		log.Println("trying to release process again")
		s.release(cmd)
		return
	}
}

func (s *service) Stop() error {

	log.Println("stopped called")

	wc := make(chan struct{})
	go func() {
		s.stop()
		wc <- struct{}{}
	}()

	select {
	case <-wc:
	case <-time.After(10 * time.Second):
		return errors.Errorf("timed out trying to stop docker tenet with id: %s", s.containerID)
	}
	return nil
}

// func (s *service) IsRunning() bool {
// 	panic("not implemented")
// }

func (s *service) DialGRPC() (*grpc.ClientConn, error) {
	c, err := s.client()
	if err != nil {
		return nil, errors.Trace(err)
	}

	dockerDialer := func(addr string, timeout time.Duration) (net.Conn, error) {
		c.Dialer.Timeout = timeout
		return c.Dialer.Dial("tcp", addr)
	}

	log.Println("dialing docker server")
	return grpc.Dial(s.ip+":"+s.port, grpc.WithDialer(dockerDialer), grpc.WithInsecure())
}
