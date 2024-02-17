package engine

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
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

func (engine *Engine) Exec() {
	res, err := engine.client.ContainerExecCreate(context.Background(), "IsaDevEnv", types.ExecConfig{
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

	conn, err := engine.client.ContainerExecAttach(context.Background(), res.ID, types.ExecStartCheck{})

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

func (engine *Engine) Attach() {
	res, err := engine.client.ContainerAttach(context.Background(),
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

// TODO might not be well-handled
func (engine *Engine) PullImage(imageName string) (<-chan string, <-chan bool, error) {
	cStatus := make(chan string)
	cDone := make(chan bool)

	rc, err := engine.client.ImagePull(context.Background(), imageName, types.ImagePullOptions{})

	if err != nil {
		return cStatus, cDone, err
	}

	go func() {
		scanner := bufio.NewScanner(rc)
		for scanner.Scan() {
			cStatus <- scanner.Text()
		}
		rc.Close()
		cDone <- true
	}()

	return cStatus, cDone, err
}

func (engine *Engine) DeleteCapsule(capsule models.Capsule) error {
	return engine.client.ContainerRemove(context.Background(), capsule.ID, types.ContainerRemoveOptions{})
}

func (engine *Engine) CreateCapsule(ccr models.CapsuleCreateRequest) (container.CreateResponse, error) {
	image := ccr.CapsuleTemplate.Namespace + "/" + ccr.CapsuleTemplate.Name

	// Pull the image
	cStatus, cDone, err := engine.PullImage(image)

	if err != nil {
		return container.CreateResponse{}, err
	}

	loop := true

	for loop {
		select {
		case status := <-cStatus:
			log.Println(status)
		case <-cDone:
			log.Println("Image pull done")
			loop = false
		}
	}

	log.Println("Creating a capsule now")

	return engine.client.ContainerCreate(context.Background(), &container.Config{
		Image:      image,
		Entrypoint: []string{"sleep", "infinity"},
		ExposedPorts: map[nat.Port]struct{}{
			"1883": {},
		},
		WorkingDir: ccr.WorkingDir,
	}, &container.HostConfig{
		Binds: []string{os.Getenv("HOME") + "/.ssh:/root/.ssh"},
	}, &network.NetworkingConfig{}, &v1.Platform{}, ccr.Name)
}

func (engine *Engine) StartEvent() (<-chan models.CapsuleEvent, <-chan error) {
	log.Println("Event has been started")

	engine.isEvent = true

	cCapsuleEvent := make(chan models.CapsuleEvent)

	cEventMessage, cError := engine.client.Events(context.Background(), types.EventsOptions{})

	go func() {
		for engine.isEvent {
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

func (engine *Engine) StopEvent() {
	engine.isEvent = false
}

func (engine *Engine) StartCapsule(id string) error {
	return engine.client.ContainerStart(context.Background(), id, types.ContainerStartOptions{})
}

func (engine *Engine) StopCapsule(id string) error {
	return engine.client.ContainerStop(context.Background(), id, container.StopOptions{})
}

func (engine *Engine) ListCapsules() ([]models.Capsule, error) {
	log.Fatalln("adsadss")
	containers, err := engine.client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	log.Fatalln(containers)
	return helpers.ConvertContainerToCapsule(containers), err
}

func CreateEngine() *Engine {
	engine := Engine{}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	engine.client = cli
	return &engine
}
