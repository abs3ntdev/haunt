package haunt

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"

	"github.com/abs3ntdev/haunt/src/config"
)

// Schema projects list
type Schema struct {
	Projects []Project `yaml:"schema" json:"schema"`
}

// Add a project if unique
func (s *Schema) Add(p Project) {
	for _, val := range s.Projects {
		if reflect.DeepEqual(val, p) {
			return
		}
	}
	s.Projects = append(s.Projects, p)
}

// Remove a project
func (s *Schema) Remove(name string) error {
	for key, val := range s.Projects {
		if name == val.Name {
			s.Projects = append(s.Projects[:key], s.Projects[key+1:]...)
			return nil
		}
	}
	return errors.New("project not found")
}

// New create a project using cli fields
func (s *Schema) New(flags config.Flags) Project {
	if flags.Name == "" {
		flags.Name = filepath.Base(Wdir())
	}
	if flags.Path == "" {
		flags.Path = Wdir()
	}
	fmt.Println(flags.Name)
	project := Project{
		Name: flags.Name,
		Path: flags.Path,
		Tools: Tools{
			Vet: Tool{
				Status: flags.Vet,
			},
			Fmt: Tool{
				Status: flags.Format,
			},
			Test: Tool{
				Status: flags.Test,
			},
			Generate: Tool{
				Status: flags.Generate,
			},
			Build: Tool{
				Status: flags.Build,
			},
			Install: Tool{
				Status: flags.Install,
			},
			Run: Tool{
				Status: flags.Run,
			},
		},
		Watcher: Watch{
			Paths:  []string{"./"},
			Ignore: []string{".git", ".haunt", "vendor"},
			Exts:   []string{"go"},
		},
	}
	return project
}

// Filter project list by names
func (s *Schema) Filter(names []string) []Project {
	result := []Project{}
	for _, item := range s.Projects {
		for _, name := range names {
			if item.Name == name {
				result = append(result, item)
			}
		}
	}
	return result
}
