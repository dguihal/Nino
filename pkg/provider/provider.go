package provider

import "github.com/dguihal/nino/pkg/docker"

type Provider interface {
	ListImages() []*docker.DockerImage
}
