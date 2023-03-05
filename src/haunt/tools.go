package haunt

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Tool info
type Tool struct {
	Args   []string `yaml:"args,omitempty" json:"args,omitempty"`
	Method string   `yaml:"method,omitempty" json:"method,omitempty"`
	Path   string   `yaml:"path,omitempty" json:"path,omitempty"`
	Dir    string   `yaml:"dir,omitempty" json:"dir,omitempty"` // wdir of the command
	Status bool     `yaml:"status,omitempty" json:"status,omitempty"`
	Output bool     `yaml:"output,omitempty" json:"output,omitempty"`
	dir    bool
	isTool bool
	cmd    []string
	name   string
	parent *Project
}

// Tools go
type Tools struct {
	Clean    Tool `yaml:"clean,omitempty" json:"clean,omitempty"`
	Vet      Tool `yaml:"vet,omitempty" json:"vet,omitempty"`
	Fmt      Tool `yaml:"fmt,omitempty" json:"fmt,omitempty"`
	Test     Tool `yaml:"test,omitempty" json:"test,omitempty"`
	Generate Tool `yaml:"generate,omitempty" json:"generate,omitempty"`
	Install  Tool `yaml:"install,omitempty" json:"install,omitempty"`
	Build    Tool `yaml:"build,omitempty" json:"build,omitempty"`
	Run      Tool `yaml:"run,omitempty" json:"run,omitempty"`
}

// Setup go tools
func (t *Tools) Setup(p *Project) {
	gocmd := "go"

	// go clean
	if t.Clean.Status {
		t.Clean.name = "Clean"
		t.Clean.isTool = true
		t.Clean.cmd = replace([]string{gocmd, "clean"}, t.Clean.Method)
		t.Clean.Args = split([]string{}, t.Clean.Args)
		msg = fmt.Sprintln(p.pname(p.Name, 1), ":", Magenta.Bold(t.Clean.name), "running. Method:", Blue.Bold(strings.Join(t.Clean.cmd, " "), " ", strings.Join(t.Clean.Args, " ")))
		log.Print(msg)
	}
	// go generate
	if t.Generate.Status {
		t.Generate.dir = true
		t.Generate.isTool = true
		t.Generate.name = "Generate"
		t.Generate.cmd = replace([]string{gocmd, "generate"}, t.Generate.Method)
		t.Generate.Args = split([]string{}, t.Generate.Args)
		msg = fmt.Sprintln(p.pname(p.Name, 1), ":", Magenta.Bold(t.Generate.name), "running. Method:", Blue.Bold(strings.Join(t.Generate.cmd, " "), " ", strings.Join(t.Generate.Args, " ")))
		log.Print(msg)
	}
	// go fmt
	if t.Fmt.Status {
		if len(t.Fmt.Args) == 0 {
			t.Fmt.Args = []string{"-s", "-w", "-e"}
		}
		t.Fmt.name = "Fmt"
		t.Fmt.isTool = true
		t.Fmt.cmd = replace([]string{"gofmt"}, t.Fmt.Method)
		t.Fmt.Args = split([]string{}, t.Fmt.Args)
		msg = fmt.Sprintln(p.pname(p.Name, 1), ":", Magenta.Bold(t.Fmt.name), "running. Method:", Blue.Bold(strings.Join(t.Fmt.cmd, " "), " ", strings.Join(t.Fmt.Args, " ")))
		log.Print(msg)
	}
	// go vet
	if t.Vet.Status {
		t.Vet.dir = true
		t.Vet.name = "Vet"
		t.Vet.isTool = true
		t.Vet.cmd = replace([]string{gocmd, "vet"}, t.Vet.Method)
		t.Vet.Args = split([]string{}, t.Vet.Args)
		msg = fmt.Sprintln(p.pname(p.Name, 1), ":", Magenta.Bold(t.Vet.name), "running. Method:", Blue.Bold(strings.Join(t.Vet.cmd, " "), " ", strings.Join(t.Vet.Args, " ")))
		log.Print(msg)
	}
	// go test
	if t.Test.Status {
		t.Test.dir = true
		t.Test.isTool = true
		t.Test.name = "Test"
		t.Test.cmd = replace([]string{gocmd, "test"}, t.Test.Method)
		t.Test.Args = split([]string{}, t.Test.Args)
		p.parent.Prefix(t.Test.name + ": Running")
		msg = fmt.Sprintln(p.pname(p.Name, 1), ":", Magenta.Bold(t.Test.name), "running. Method:", Blue.Bold(strings.Join(t.Test.cmd, " "), " ", strings.Join(t.Test.Args, " ")))
		log.Print(msg)
	}
	// go install
	t.Install.name = "Install"
	t.Install.cmd = replace([]string{gocmd, "install"}, t.Install.Method)
	t.Install.Args = split([]string{}, t.Install.Args)
	// go build
	if t.Build.Status {
		t.Build.name = "Build"
		t.Build.cmd = replace([]string{gocmd, "build"}, t.Build.Method)
		t.Build.Args = split([]string{}, t.Build.Args)
	}
}

// Exec a go tool
func (t *Tool) Exec(path string, stop <-chan bool) (response Response) {
	if t.dir {
		if filepath.Ext(path) != "" {
			path = filepath.Dir(path)
		}
		// check if there is at least one go file
		matched := false
		files, _ := os.ReadDir(path)
		for _, f := range files {
			matched, _ = filepath.Match("*.go", f.Name())
			if matched {
				break
			}
		}
		if !matched {
			return
		}
	} else if !strings.HasSuffix(path, ".go") {
		return
	}
	args := t.Args
	if strings.HasSuffix(path, ".go") {
		args = append(args, path)
		path = filepath.Dir(path)
	}
	if s := ext(path); s == "" || s == "go" {
		if t.parent.parent.Settings.Recovery.Tools {
			log.Println("Tool:", t.name, path, args)
		}
		var out, stderr bytes.Buffer
		done := make(chan error)
		args = append(t.cmd, args...)
		cmd := exec.Command(args[0], args[1:]...)
		if t.Dir != "" {
			cmd.Dir, _ = filepath.Abs(t.Dir)
		} else {
			cmd.Dir = path
		}
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		// Start command
		err := cmd.Start()
		if err != nil {
			response.Name = t.name
			response.Err = err
			return
		}
		go func() { done <- cmd.Wait() }()
		// Wait a result
		select {
		case <-stop:
			// Stop running command
			if err := cmd.Process.Kill(); err != nil {
				response.Err = err
			}
		case err := <-done:
			// Command completed
			response.Name = t.name
			if err != nil {
				response.Err = errors.New(stderr.String() + out.String() + err.Error())
			} else {
				if t.Output {
					response.Out = out.String()
				}
			}
		}
	}
	return
}

// Compile is used for build and install
func (t *Tool) Compile(path string, stop <-chan bool) (response Response) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	done := make(chan error)
	args := append(t.cmd, t.Args...)
	cmd := exec.Command(args[0], args[1:]...)
	if t.Dir != "" {
		cmd.Dir, _ = filepath.Abs(t.Dir)
	} else {
		cmd.Dir = path
	}
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	// Start command
	err := cmd.Start()
	if err != nil {
		log.Println(log.Prefix(), err.Error())
		if cmd.Process != nil {
			err := cmd.Process.Kill()
			if err != nil {
				log.Println(log.Prefix(), err.Error())
			}
		}
		os.Exit(0)
	}
	go func() {
		done <- cmd.Wait()
	}()
	// Wait a result
	response.Name = t.name
	select {
	case <-stop:
		// Stop running command
		err := cmd.Process.Kill()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	case err := <-done:
		// Command completed
		if err != nil {
			response.Err = errors.New(stderr.String() + err.Error())
		}
	}
	return
}
