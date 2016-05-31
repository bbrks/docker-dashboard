package main

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/fsouza/go-dockerclient"
)

var (
	t *template.Template
	c *docker.Client
)

type page struct {
	C         []docker.APIContainers
	Container docker.Container
	CID, H, L string
}

func main() {
	var err error
	t, err = template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}

	c, err = docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/logs/", logHandler)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func getContainers() []docker.APIContainers {
	containers, err := c.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		panic(err)
	}

	return containers
}

func getLogs(cID string) string {
	var buf bytes.Buffer
	c.Logs(docker.LogsOptions{
		Container:    cID,
		OutputStream: &buf,
		Stdout:       true,
		Stderr:       true,
		RawTerminal:  true,
		Timestamps:   true,
		Tail:         "50",
	})
	return buf.String()
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	cID := strings.Split(r.URL.Path, "/")[2]
	p := page{
		L: getLogs(cID),
		H: r.Host,
	}
	err := t.Execute(w, p)
	if err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := page{C: getContainers(), H: r.Host}
	err := t.Execute(w, p)
	if err != nil {
		panic(err)
	}
}
