package docker

import (
	"log/slog"
	"strings"

	semver "github.com/Masterminds/semver/v3"
)

type DockerImage struct {
	Image   string
	Version string
}

func NewDockerImageFromString(input string) (DockerImage, bool) {

	if len(input) == 0 {
		return DockerImage{}, false
	}

	x := strings.Split(input, ":")
	dockerImage := DockerImage{Image: x[0], Version: ""}
	if len(x) > 1 {
		dockerImage.Version = x[1]
	}

	return dockerImage, true
}

// Keep only semver compatible tagged images
func SemverFilteredDockerImages(images []*DockerImage) []*DockerImage {

	var filteredImages []*DockerImage

	for _, image := range images {

		slog.Debug("semverFilteredDockerImages: input task image", "Image", image.Image, "Version", image.Version)

		if image.Version == "" || image.Version == "latest" {
			continue
		}

		if _, err := semver.NewVersion(image.Version); err != nil {
			continue
		}

		slog.Info("Semver compatible task image found:", "Image", image.Image, "Version", image.Version)

		image := DockerImage{Image: image.Image, Version: image.Version}

		filteredImages = append(filteredImages, &image)
	}
	return filteredImages
}

// Deduplicate list of docker images
func DedupFilteredDockerImages(images []*DockerImage) []*DockerImage {

	var filteredImages []*DockerImage

	for _, image := range images {

		found := false

		for index, filImage := range filteredImages {
			if image.Image == filImage.Image {
				if image.Version < filImage.Version {
					filteredImages[index] = image
				}
				found = true
				break
			}
		}

		if !found {
			filteredImages = append(filteredImages, image)
		}
	}

	return filteredImages
}
