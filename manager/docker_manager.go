package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const CONTAINER_NAME string = "docker_manager"

type Subscribers interface {
	NewContainer(dockerPort int)
}

type DockerManager struct {
	cli             *client.Client
	subscribers     []Subscribers
	image           string
	appPort         int
	oldContainerId  *string
	overlapDuration time.Duration
}

func NewDockerManager(cli *client.Client, image string, appPort int,
	overlapDuration time.Duration) *DockerManager {

	return &DockerManager{
		cli:             cli,
		appPort:         appPort,
		image:           image,
		overlapDuration: overlapDuration,
	}
}

func (dm *DockerManager) AddSubscriber(s Subscribers) {
	dm.subscribers = append(dm.subscribers, s)
}

func getFreePort() (int, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil

}

func (dm *DockerManager) runNewContainer(newImageID string) {
	ctx := context.Background()
	defer ctx.Done()

	port, _ := nat.NewPort("tcp", strconv.Itoa(dm.appPort))
	hostPort, err := getFreePort()
	if err != nil {
		log.Println("Can't find a free port", err.Error())
	}

	container, err := dm.cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: newImageID,
			ExposedPorts: nat.PortSet{
				port: struct{}{},
			},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				port: []nat.PortBinding{
					{
						HostPort: strconv.Itoa(hostPort),
					},
				},
			},
			AutoRemove: true,
		},
		&network.NetworkingConfig{},
		fmt.Sprintf("%s_%d", CONTAINER_NAME, hostPort),
	)
	if err != nil {
		log.Println("Couldn't create the container", err.Error())
		return
	}
	log.Printf("Container created with ID: %s; to listen on port %d\n", container.ID, hostPort)
	err = dm.cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})

	if err != nil {
		log.Println("Error starting the container", err.Error())
		return
	}

	for _, s := range dm.subscribers {
		s.NewContainer(hostPort)
	}

	// Do nothing for desired duration
	time.Sleep(dm.overlapDuration)

	// Remove old containers
	list, err := dm.cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Println("Couldn't list containers", err.Error())
		return
	}

	for _, c := range list {
		delete := false
		var n string
		for _, n = range c.Names {
			if c.ID != container.ID && strings.Contains(n, CONTAINER_NAME) {
				delete = true
				break
			}
		}
		if delete {
			log.Printf("Removing container id: %s (name: %s)", c.ID, n)
			d := time.Second * 5
			err := dm.cli.ContainerStop(ctx, c.ID, &d)
			if err != nil {
				log.Println("Cound not stop container", err.Error())
			}
		}
	}
}

func (dm *DockerManager) Run(updates <-chan string, wg *sync.WaitGroup, exit <-chan bool) {
	for {
		select {
		case _ = <-exit:
			log.Println("Received exit request; docker manager")
			wg.Done()
			return
		case update := <-updates:
			dm.runNewContainer(update)
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
