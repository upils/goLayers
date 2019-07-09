package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
)

var layerMap map[string][]string


func displayNamespaces (client *containerd.Client) error {
	nsStore := client.NamespaceService()
	labels, err := nsStore.List(context.Background())
	if err != nil {
		return err
	}

	for _, label := range labels {
		fmt.Printf("%s\n", label)
	}
	return nil
}

func main() {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
			panic(err)
	}
	defer client.Close()

	displayNamespaces(client)

	ctx := namespaces.WithNamespace(context.Background(), "moby")

	images, err := client.ListImages(ctx)
	if err != nil {
		panic(err)
	}

	layerMap = make(map[string][]string)

	for _, image := range images {
		layers, err := image.RootFS(ctx)
		
		if err != nil {
			panic(err)
		}
		for _, layer := range layers {
			var layerString = layer.String()
			fmt.Printf(layerString)
			layerMap[layerString] = append(layerMap[layerString], image.Name())
		}
	}

	fmt.Printf("Layer ID -> Images\n")
	for layer, imageList := range layerMap {
		fmt.Printf("%s -> %s\n", layer[7:19], strings.Join(imageList , ", "))
    }
}