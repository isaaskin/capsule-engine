package engine

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/isaaskin/capsule-engine/helpers"
	"github.com/isaaskin/capsule-engine/models"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Binding[T int | string] struct {
	container T
	local     T
}

type MyDevEnv struct {
	Name    string
	Image   string
	Ports   []Binding[int]
	Volumes []Binding[string]
	IsGUI   bool
}

func (m *MyDevEnv) build() *container.Config {
	return &container.Config{
		Hostname: m.Name,
		Image:    m.Image,
	}
}

type Engine struct {
	client  *client.Client
	isEvent bool
}

func (Engine *Engine) Exec() {
	res, err := Engine.client.ContainerExecCreate(context.Background(), "IsaDevEnv", types.ExecConfig{
		Tty:          true,
		AttachStdin:  true,
		AttachStderr: false,
		AttachStdout: true,
		Detach:       true,
		WorkingDir:   "",
		Cmd:          []string{"/bin/bash"},
	})

	if err != nil {
		fmt.Println("error here....")
		log.Fatalf(err.Error())
	}

	conn, err := Engine.client.ContainerExecAttach(context.Background(), res.ID, types.ExecStartCheck{})

	if err != nil {
		fmt.Println("error here....")
		log.Fatalf(err.Error())
	}

	go func() {
		scanner := bufio.NewScanner(conn.Reader)
		for scanner.Scan() {
			fmt.Println("-> " + scanner.Text())
			fmt.Println("-----------------")
		}
	}()

	if err != nil {
		log.Fatalf(err.Error())
	}

	for {
		time.Sleep(1 * time.Second)
		_, err = conn.Conn.Write([]byte("echo 1\r"))
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func (Engine *Engine) Attach() {
	res, err := Engine.client.ContainerAttach(context.Background(),
		"IsaDevEnv",
		types.ContainerAttachOptions{
			Stream: true,
			Stdin:  true,
			Stdout: true,
			Stderr: true,
			Logs:   false,
		})

	if err != nil {
		log.Fatalf(err.Error())
	}

	go func() {
		scanner := bufio.NewScanner(res.Reader)
		for scanner.Scan() {
			fmt.Println("-> " + scanner.Text())
		}
	}()

	_, err = res.Conn.Write([]byte("/bin/bash\n"))

	if err != nil {
		log.Fatalf(err.Error())
	}

	go func() {
		for {
			_, err = res.Conn.Write([]byte("echo 123\n"))
			if err != nil {
				log.Fatalf(err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (Engine *Engine) CreateCapsule(m *MyDevEnv) (container.CreateResponse, error) {
	rc, err := Engine.client.ImagePull(context.Background(), m.Image, types.ImagePullOptions{})
	if err != nil {
		return container.CreateResponse{}, err
	}

	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return Engine.client.ContainerCreate(context.Background(), &container.Config{
		Image:      m.Image,
		Entrypoint: []string{"sleep", "infinity"},
		ExposedPorts: map[nat.Port]struct{}{
			"1883": {},
		},
	}, &container.HostConfig{}, &network.NetworkingConfig{}, &v1.Platform{}, m.Name)
}

func (Engine *Engine) StartEvent() (<-chan models.CapsuleEvent, <-chan error) {
	log.Println("Event has been started")

	Engine.isEvent = true

	cCapsuleEvent := make(chan models.CapsuleEvent)

	cEventMessage, cError := Engine.client.Events(context.Background(), types.EventsOptions{})

	go func() {
		for Engine.isEvent {
			select {
			case eventMessage := <-cEventMessage:
				cCapsuleEvent <- models.CapsuleEvent{
					Target: eventMessage.Actor.Attributes["name"],
					Action: eventMessage.Action,
				}
			}
		}
		log.Println("Event has been stopped")
	}()

	return cCapsuleEvent, cError
}

func (Engine *Engine) StopEvent() {
	Engine.isEvent = false
}

func (Engine *Engine) StartCapsule(id string) error {
	return Engine.client.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

func (Engine *Engine) StopCapsule(id string) error {
	return Engine.client.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (Engine *Engine) ListCapsules() ([]models.Capsule, error) {
	containers, err := Engine.client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	return helpers.ConvertContainerToCapsule(containers), err
}

func CreateEngine() *Engine {
	Engine := Engine{}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	Engine.client = cli
	return &Engine
}
