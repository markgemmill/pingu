package pkg

import (
	"embed"
	"errors"
	"fmt"
	iofs "io/fs"
	"path"
)

var (
	//go:embed templates/*
	emfs embed.FS
)

type Resources struct {
	directory string
	templates []iofs.DirEntry
}

func (r *Resources) Initialize() {
	templates, _ := iofs.ReadDir(emfs, r.directory)
	r.templates = templates
}

func (r *Resources) ReadText(name string) (string, error) {
	for _, template := range r.templates {
		//fmt.Printf("Resource: %s\n", template.Name())
		if template.Name() == name && !template.Type().IsDir() {
			text, err := iofs.ReadFile(emfs, path.Join(r.directory, template.Name()))
			if err != nil {
				return "", err
			}
			return string(text), nil
		}
	}
	return "", errors.New(fmt.Sprintf("No template: %s", name))
}

func NewResources(directory string) *Resources {
	res := Resources{directory: directory}
	res.Initialize()
	return &res
}
