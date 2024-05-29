package nino

import (
	"context"
	"log/slog"
	"os"

	"github.com/dguihal/nino/pkg/docker"
	"github.com/dguihal/nino/pkg/provider"
	"github.com/dguihal/nino/pkg/provider/nomad"
)

const nomadEnvKey = "NOMAD_ADDR"

func Check(ctx context.Context) {

	var images []*docker.DockerImage

	providers := buildProviders()

	slog.Info("Step 1: Getting list of image from providers ...")
	for _, p := range providers {
		images = append(images, p.ListImages()...)
	}
	slog.Debug("Pre filtered image count to check", "count", len(images))

	slog.Info("Step 2: Filtering gathered image list ...")
	images = filterDockerImages(images)

	slog.Info("Post filtered image count to check", "count", len(images))

	slog.Info("Step 3: Looking for new image tags ...")
	for _, image := range images {

		slog.Info("Looking for new image tags for image", "image", image.Image)
		slog.Info("Reference version", "version", image.Version)
		docker.CheckForNewVersion(image, ctx)
	}
}

// Keep only semver compatible tagged images
func filterDockerImages(images []*docker.DockerImage) []*docker.DockerImage {

	filteredImages := docker.SemverFilteredDockerImages(images)

	filteredImages = docker.DedupFilteredDockerImages(filteredImages)

	return filteredImages
}

func buildProviders() []provider.Provider {

	var providers []provider.Provider

	if len(os.Getenv(nomadEnvKey)) > 0 {
		slog.Info("Adding new nomad provider", "address", os.Getenv(nomadEnvKey))
		nomadProvider := nomad.New()
		providers = append(providers, nomadProvider)
	}

	return providers
}
