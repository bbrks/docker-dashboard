package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"golang.org/x/net/context"
)

var t *template.Template

type page struct {
	C []types.Container
	H string
}

func main() {
	var err error
	t, err = template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func getContainers() []types.Container {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClient("unix:///var/run/docker.sock", "v1.22", nil, defaultHeaders)
	if err != nil {
		panic(err)
	}

	options := types.ContainerListOptions{All: true}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		c.Image = strings.TrimPrefix(c.Image, os.Getenv("TRIM_IMAGE_PREFIX"))
	}

	return containers
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := page{C: getContainers(), H: r.Host}
	err := t.Execute(w, p)
	if err != nil {
		panic(err)
	}
}
