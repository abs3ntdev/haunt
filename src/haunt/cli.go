package haunt

import (
	"errors"
	"fmt"
	"go/build"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	// HPrefix tool name
	HPrefix = "haunt"
	// HExt file extension
	HExt = ".yaml"
	// HFile config file name
	HFile = "." + HPrefix + HExt
	// HExtWin windows extension
	HExtWin = ".exe"
)

type (
	// LogWriter used for all log
	LogWriter struct{}

	// Haunt main struct
	Haunt struct {
		Settings Settings `yaml:"settings" json:"settings"`
		Server   Server   `yaml:"server,omitempty" json:"server,omitempty"`
		Schema   `yaml:",inline" json:",inline"`
		Sync     chan string `yaml:"-" json:"-"`
		Err      Func        `yaml:"-" json:"-"`
		After    Func        `yaml:"-"  json:"-"`
		Before   Func        `yaml:"-"  json:"-"`
		Change   Func        `yaml:"-"  json:"-"`
		Reload   Func        `yaml:"-"  json:"-"`
	}

	// Context is used as argument for func
	Context struct {
		Path    string
		Project *Project
		Stop    <-chan bool
		Watcher FileWatcher
		Event   fsnotify.Event
	}

	// Func is used instead haunt func
	Func func(Context)
)

func NewHaunt() *Haunt {
	return &Haunt{
		Sync: make(chan string),
	}
}

// init check
func init() {
	// custom log
	log.SetFlags(0)
	log.SetOutput(LogWriter{})
	if build.Default.GOPATH == "" {
		log.Fatal("$GOPATH isn't set properly")
	}
	path := filepath.SplitList(build.Default.GOPATH)
	if err := os.Setenv("GOBIN", filepath.Join(path[len(path)-1], "bin")); err != nil {
		log.Fatal(err)
	}
}

func (h *Haunt) SetDefaults() {
	h.Server = Server{Parent: h, Status: true, Open: false, Port: Port}
	h.Settings.FileLimit = 0
	h.Settings.Legacy.Interval = 100 * time.Millisecond
	h.Settings.Legacy.Force = false
	h.Settings.Errors = Resource{Name: FileErr, Status: false}
	h.Settings.Errors = Resource{Name: FileOut, Status: false}
	h.Settings.Errors = Resource{Name: FileLog, Status: false}
	if _, err := os.Stat("main.go"); err == nil {
		log.Println(h.Prefix(Green.Bold("Adding: " + filepath.Base(Wdir()))))
		h.Projects = append(h.Projects, Project{
			Name: filepath.Base(Wdir()),
			Path: Wdir(),
			Tools: Tools{
				Install: Tool{
					Status: true,
				},
				Run: Tool{
					Status: true,
				},
			},
			Watcher: Watch{
				Exts:  []string{"go"},
				Paths: []string{"./"},
			},
		})
	} else {
		log.Println(h.Prefix(Magenta.Bold("Skipping: " + filepath.Base(Wdir()) + " no main.go file in root")))
	}
	subDirs, err := os.ReadDir("cmd")
	if err != nil {
		log.Println(h.Prefix("cmd directory not found, skipping"))
		return
	}
	for _, dir := range subDirs {
		if dir.IsDir() {
			log.Println(h.Prefix(Green.Bold("Adding: " + dir.Name())))
			h.Projects = append(h.Projects, Project{
				Name: dir.Name(),
				Path: "cmd/" + dir.Name(),
				Tools: Tools{
					Install: Tool{
						Status: true,
					},
					Run: Tool{
						Status: true,
					},
				},
				Watcher: Watch{
					Exts:  []string{"go"},
					Paths: []string{"cmd/" + dir.Name()},
				},
			})
		} else {
			log.Println(h.Prefix(Magenta.Bold("Skipping: " + dir.Name() + " not a directory")))
		}
	}
}

// Stop haunt workflow
func (h *Haunt) Stop() error {
	for k := range h.Projects {
		if h.Schema.Projects[k].exit != nil {
			close(h.Schema.Projects[k].exit)
		}
	}
	return nil
}

// Run haunt workflow
func (h *Haunt) Run() error {
	if len(h.Projects) > 0 {
		var wg sync.WaitGroup
		wg.Add(len(h.Projects))
		for k := range h.Projects {
			h.Schema.Projects[k].exit = make(chan os.Signal, 1)
			signal.Notify(h.Schema.Projects[k].exit, os.Interrupt)
			h.Schema.Projects[k].parent = h
			go h.Schema.Projects[k].Watch(&wg)
		}
		wg.Wait()
	} else {
		return errors.New("there are no projects")
	}
	return nil
}

// Prefix a given string with tool name
func (h *Haunt) Prefix(input string) string {
	if len(input) > 0 {
		return fmt.Sprint(Yellow.Bold("["), strings.ToUpper(HPrefix), Yellow.Bold("]"), ": ", input)
	}
	return input
}

// Rewrite the layout of the log timestamp
func (w LogWriter) Write(bytes []byte) (int, error) {
	if len(bytes) > 0 {
		return fmt.Fprint(Output, Yellow.Regular("["), time.Now().Format("15:04:05"), Yellow.Regular("]"), string(bytes))
	}
	return 0, nil
}
