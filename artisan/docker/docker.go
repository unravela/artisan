package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/unravela/artisan/api"
	"os"
)

// ImagePrefix is used in all docker images generated from source. If
// env is named as "node", the image will be "artsn-node:latest"
const ImagePrefix = "artsn-"

func pullImage(docker *client.Client, image string) (api.EnvironmentID, error) {

	if image == "" {
		return "", errors.New("no image is present")
	}

	// We need canonical name of the image. That means not only something:latest but
	// also registry need to be part of the string. For now,
	// we are using docker.io/library as default
	img := "docker.io/library/" + image

	ctx := context.Background()
	res, err := docker.ImagePull(ctx, img, types.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("cannot pull image %w", err)
	}

	// print the image build result to output
	scanner := bufio.NewScanner(res)
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Printf("\033[2K\r    > %s", txt)
	}

	imageID := getImageID(docker, image)
	envID := api.EnvironmentID(imageID)
	return envID, nil
}

// this function builds docker image for given source dir. The image is build
// with given name
func buildImage(docker *client.Client, name string, srcDir string) (api.EnvironmentID, error) {
	ctx := context.Background()
	// create tar
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return "", fmt.Errorf("invalid path to source dir")
	}

	bctx := createBuildContext(buildContextOptions{
		root:     srcDir,
		includes: []string{},
	})

	// build image
	tags := []string{envToTag(name)}
	res, err := docker.ImageBuild(ctx, bctx, types.ImageBuildOptions{
		Dockerfile: "./Dockerfile",
		NoCache:    true,
		Tags:       tags,
	})

	if err != nil {
		return "", fmt.Errorf("Cannot build env %w", err)
	}

	// print the image build result to output
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Printf("\033[2K\r    > %s", txt)
	}

	imageID := getImageID(docker, envToTag(name))
	envID := api.EnvironmentID(imageID)
	return envID, nil
}

// function returns you docker ID for given image name
// This is needed e.g. when you want to perform task and task is
// running in forge.
//
// If function returns you empty string, that means there is no
// docker image present in system for this forge and forge need to
// be build.
func getImageID(docker *client.Client, image string) string {
	ctx := context.Background()
	images, _ := docker.ImageList(ctx, types.ImageListOptions{})
	for _, img := range images {
		tags := img.RepoTags

		idx := len(tags) - 1
		if idx < 0 {
			idx = 0
		}

		for _, tag := range tags {
			if tag == image {
				return img.ID
			}
		}
	}

	return ""
}

// This function transform environment name e.g. 'jdk8' to
// docker tag 'dlvin-jdk8:latest'
func  envToTag(name string) string {
	return ImagePrefix + name + ":latest"
}
