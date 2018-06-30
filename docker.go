package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Containers map[string]int

func (c *Containers) Add(name string) (err error) {
	if strings.HasPrefix(name, "/") {
		name = strings.Replace(name, "/", "", 1)
	}

	_, ok := (*c)[name]
	if ok {
		return fmt.Errorf("Container %q already registered", name)
	}

	(*c)[name] = 1

	return
}

type DockerClient interface {
	ContainerList(context.Context, types.ContainerListOptions) ([]types.Container, error)
}

type Docker struct {
	Client DockerClient
}

func NewDocker() (d Docker, err error) {
	d.Client, err = client.NewEnvClient()

	return
}

func (d Docker) Containers() (c Containers, err error) {
	containers, err := d.Client.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return
	}

	c = make(Containers)

	for _, container := range containers {
		for _, n := range container.Names {
			err = c.Add(n)
			if err != nil {
				return
			}
		}
	}

	return
}
