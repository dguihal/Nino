package docker

import (
	"context"
	"log/slog"
	"os"

	semver "github.com/Masterminds/semver/v3"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func CheckForNewVersion(dockerImage *DockerImage, ctx context.Context) {

	var found []semver.Version

	checkConstraint, err := semver.NewConstraint("> " + dockerImage.Version)
	if err != nil {
		slog.Error("Invalid docker version", "image", dockerImage.Image, "version", dockerImage.Version, "error", err)
		os.Exit(1)
	}

	repo, err := name.NewRepository(dockerImage.Image)
	if err != nil {
		slog.Error("Error building repository object", "repo", dockerImage.Image, "error", err)
		os.Exit(1)
	}

	puller, err := remote.NewPuller()
	if err != nil {
		slog.Error("Error creating puller", err)
		os.Exit(1)
	}

	lister, _ := puller.Lister(ctx, repo)

	for lister.HasNext() {
		tags, _ := lister.Next(ctx)
		for _, tag := range tags.Tags {

			v, err := semver.NewVersion(tag)

			if err != nil || v.Prerelease() != "" {
				continue
			}

			if checkConstraint.Check(v) {
				alreadyFound := false

				for _, foundV := range found {
					if foundV.Equal(v) {
						alreadyFound = true
						break
					}
				}

				// Needed to avoid vX.X.X and X.X.X shown independently (and having duplicate output)
				if !alreadyFound {
					found = append(found, *v)
					slog.Info("Found value", "Path", dockerImage.Image, "Ref", refV.String(), "Tag", tag)
				}
			}

		}
	}

	if len(found) == 0 {
		slog.Info("No new tag found")
	}
}
