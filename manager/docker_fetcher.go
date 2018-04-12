package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const HTTP_TIMEOUT time.Duration = time.Second

type DockerFetcher struct {
	cli          *client.Client
	tickInterval time.Duration
	image        string
	ch           chan string
	latestID     string
}

func (c *DockerFetcher) Updates() <-chan string {
	return c.ch
}

func (c *DockerFetcher) fetchLatest() {
	ctx := context.Background()
	defer ctx.Done()

	refStr := fmt.Sprintf("%s:latest", c.image)
	o, err := c.cli.ImagePull(ctx, refStr, types.ImagePullOptions{})
	if o != nil {
		defer o.Close()
	}
	if err != nil {
		log.Println("Couldn't pull the latest image", err.Error())
		return
	}
	inspect, _, err := c.cli.ImageInspectWithRaw(ctx, refStr)
	if err != nil {
		log.Println("Couldn't pull the latest image", err.Error())
		return
	}
	if c.latestID != inspect.ID {
		log.Println("New image found", inspect.ID)
		c.ch <- inspect.ID
		c.latestID = inspect.ID
	}
}

func NewDockerFetcher(cli *client.Client, image string, tickInterval time.Duration) *DockerFetcher {
	return &DockerFetcher{
		cli:          cli,
		tickInterval: tickInterval,
		image:        image,
		ch:           make(chan string, 1),
	}
}

func (c *DockerFetcher) Run(wg *sync.WaitGroup, exit <-chan bool) {
	c.fetchLatest()
	tick := time.Tick(c.tickInterval)
	for {
		select {
		case _ = <-exit:
			log.Println("Received exit request; docker fetcher")
			wg.Done()
			return
		case _ = <-tick:
			c.fetchLatest()
		default:
			time.Sleep(time.Millisecond * 50)
		}
	}
}
