package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

var layerMap map[string][]string

func main() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	layerMap = make(map[string][]string)

	for _, image := range images {
		imageInspect, _, err := cli.ImageInspectWithRaw(context.Background(), image.ID)
		if err != nil {
			panic(err)
		}
		for _, layer := range imageInspect.RootFS.Layers {
			layerMap[layer] = append(layerMap[layer], strings.Join(image.RepoTags , ", "))
		}
	}

	fmt.Printf("Layer ID -> Images\n")
	for layer, imageList := range layerMap {
		fmt.Printf("%s -> %s\n", layer[7:19], strings.Join(imageList , ", "))
    }
}