package nomad

import (
	"log/slog"
	"os"

	"github.com/dguihal/nino/pkg/docker"
	"github.com/hashicorp/nomad/api"
)

type NomadProvider struct {
}

func New() *NomadProvider {
	return &NomadProvider{}
}

func (p *NomadProvider) ListImages() []*docker.DockerImage {

	var dockerImages []*docker.DockerImage

	config := &api.Config{
		//		Address: "http://nomad.lan.ototu.me",
	}

	client, err := api.NewClient(config)
	if err != nil {
		slog.Error("Error creating nomad client", slog.Any("error", err))
		os.Exit(1)
	}

	tasks := p.getDockerTasks(client)

	for _, task := range tasks {
		dockerImage := p.dockerImageFromNomadTask(task)
		if dockerImage != nil {
			dockerImages = append(dockerImages, dockerImage)
		}
	}

	return dockerImages
}

func (p *NomadProvider) getDockerTasks(client *api.Client) []*api.Task {

	var tasks []*api.Task

	jobs, _, err := client.Jobs().List(nil)

	if err != nil {
		slog.Error("Error getting list of nomad jobs", slog.Any("error", err))
		os.Exit(1)
	}

	for index := range jobs {
		slog.Debug("job", "name", jobs[index].Name, "id", jobs[index].ID)

		jobinfos, _, err := client.Jobs().Info(jobs[index].ID, nil)

		if err != nil {
			slog.Error("Error getting job Info", slog.Any("error", err))
			os.Exit(1)
		}

		for k, v := range jobinfos.Meta {
			slog.Info("meta", k, v)
		}

		for _, taskGroup := range jobinfos.TaskGroups {

			slog.Debug("Task groups", "Name", *taskGroup.Name)

			for _, task := range taskGroup.Tasks {

				if task.Driver == "docker" {
					tasks = append(tasks, task)
				}
			}
		}
	}

	return tasks
}

func (p *NomadProvider) dockerImageFromNomadTask(task *api.Task) *docker.DockerImage {
	imageName, ok := task.Config["image"]
	if !ok {
		return nil
	}

	dockerImage, ok := docker.NewDockerImageFromString(imageName.(string))
	if !ok {
		return nil
	}

	return &dockerImage
}
