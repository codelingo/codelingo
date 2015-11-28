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

	// start a new container
	opts := goDocker.CreateContainerOptions{
		Config: &goDocker.Config{
			Image:       s.image,
			PortSpecs:   []string{"8000/tcp"},
			AttachStdin: true,
			Tty:         true,
		},
		HostConfig: &goDocker.HostConfig{
			Binds: []string{pwd + ":/source:ro"},
			PortBindings: map[goDocker.Port][]goDocker.PortBinding{
				"8000/tcp": []goDocker.PortBinding{{
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

	log.Printf("waiting for ports to bind")
	for container.NetworkSettings.Ports["8000/tcp"] == nil {
		time.Sleep(1 * time.Microsecond)
		container, err = c.InspectContainer(container.ID)
		if err != nil {
			return errors.Trace(err)
		}
	}

	log.Printf("%#v", container.NetworkSettings.Ports["8000/tcp"])

	ports := container.NetworkSettings.Ports["8000/tcp"]
	s.ip = ports[0].HostIP
	s.port = ports[0].HostPort

	log.Print("got to end of docker service.Start, no errors")
	return nil
}

func (s *service) client() (*goDocker.Client, error) {
	if s.dockerClient == nil {
		// TODO(waigani) get endpoint from ~/.lingo/config.toml
		endpoint := "unix:///var/run/docker.sock"

		dClient, err := goDocker.NewClient(endpoint)
		if err != nil {
			return nil, err
		}
		s.dockerClient = dClient
	}

	return s.dockerClient, nil
}

func (s *service) Stop() error {

	// Use exec so it's quick
	log.Println("stopped called")
	cmd := exec.Command("docker", "rm", "-f", s.containerID)
	cmd.Start()
	cmd.Process.Release()

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
