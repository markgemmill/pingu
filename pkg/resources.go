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

// Resources represents a resource directory.
type Resources struct {
	directory string
	templates []iofs.DirEntry
}

func (r *Resources) Initialize() {
	templates, _ := iofs.ReadDir(emfs, r.directory)
	r.templates = templates
}

// ReadText opens and reads the named resource from the directory as text.
func (r *Resources) ReadText(name string) (string, error) {
	for _, template := range r.templates {
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

// NewResources creates a new Resource object.
func NewResources(directory string) *Resources {
	res := Resources{directory: directory}
	res.Initialize()
	return &res
}
