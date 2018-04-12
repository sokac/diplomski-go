package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/client"
)

func signalHandler() <-chan bool {
	ch := make(chan bool)
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-s // consume channel
		log.Println("Received signal", sig)
		close(ch)
	}()
	return ch
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s path/to/config.json\n", os.Args[0])
	}
	config, err := loadConfig(os.Args[1])
	if err != nil {
		log.Fatalf("Couldn't load config %s; error: \n", os.Args[1], err.Error())
	}
	image := config.DockerImage
	if image == "" {
		log.Fatal("DOCKER_IMAGE not defined")
	}

	dockerCli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal("Can't create Docker client")
	}

	ch := signalHandler()

	c := NewDockerFetcher(dockerCli, image, time.Second*30)

	t := time.Second * time.Duration(config.VersionOverlapDuration)
	dm := NewDockerManager(dockerCli, image, 8888, t)

	dm.AddSubscriber(NewNginxConfiguration(config.NginxPIDFile, config.NginxConfiguration))

	wg := sync.WaitGroup{}
	wg.Add(2)
	go dm.Run(c.Updates(), &wg, ch)

	go c.Run(&wg, ch)

	wg.Wait()
	log.Println("Exiting Manager")
}
