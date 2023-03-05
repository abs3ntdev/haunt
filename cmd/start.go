package cmd

import (
	"log"
	"strings"

	"github.com/abs3ntdev/haunt/src/config"
	"github.com/abs3ntdev/haunt/src/haunt"
	"github.com/spf13/cobra"
)

var startConfig config.Flags

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start haunt on a given path, generates a config file if one does not already exist",
	Args:  cobra.MatchAll(cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getProjectNamesToStart(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: start,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&startConfig.Path, "path", "p", "", "Project base path")
	startCmd.Flags().BoolVarP(&startConfig.Format, "fmt", "f", false, "Enable go fmt")
	startCmd.Flags().BoolVarP(&startConfig.Vet, "vet", "v", false, "Enable go vet")
	startCmd.Flags().BoolVarP(&startConfig.Test, "test", "t", false, "Enable go test")
	startCmd.Flags().BoolVarP(&startConfig.Generate, "generate", "g", false, "Enable go generate")
	startCmd.Flags().BoolVarP(&startConfig.Server, "server", "s", false, "Start server")
	startCmd.Flags().BoolVarP(&startConfig.Open, "open", "o", false, "Open into the default browser")
	startCmd.Flags().BoolVarP(&startConfig.Install, "install", "i", true, "Enable go install")
	startCmd.Flags().BoolVarP(&startConfig.Build, "build", "b", true, "Enable go build")
	startCmd.Flags().BoolVarP(&startConfig.Run, "run", "r", true, "Enable go run")
	startCmd.Flags().BoolVarP(&startConfig.Legacy, "legacy", "l", false, "Legacy watch by polling instead fsnotify")
	startCmd.Flags().BoolVarP(&startConfig.NoConfig, "no-config", "c", false, "Ignore existing config and doesn't create a new one")
}

func getProjectNamesToStart(input string) []string {
	r := haunt.NewHaunt()
	// read a config if exist
	err := r.Settings.Read(&r)
	if err != nil {
		return []string{}
	}
	names := []string{}
	for _, project := range r.Projects {
		if strings.HasPrefix(project.Name, input) {
			names = append(names, project.Name)
		}
	}
	return names
}

// Start haunt workflow
func start(cmd *cobra.Command, args []string) (err error) {
	r := haunt.NewHaunt()
	// set legacy watcher
	if startConfig.Legacy {
		r.Settings.Legacy.Set(startConfig.Legacy, 1)
	}
	// set server
	if startConfig.Server {
		r.Server.Set(startConfig.Server, startConfig.Open, haunt.Port, haunt.Host)
	}

	// check no-config and read
	if !startConfig.NoConfig {
		// read a config if exist
		err = r.Settings.Read(&r)
		if err != nil {
			return err
		}
		if len(args) >= 1 {
			// filter by name flag if exist
			r.Projects = r.Filter(args)
			if len(r.Projects) == 0 {
				log.Println(r.Prefix("Project not found, exiting"))
				return
			}
			startConfig.Name = args[0]
		}
		// increase file limit
		if r.Settings.FileLimit != 0 {
			if err = r.Settings.Flimit(); err != nil {
				return err
			}
		}

	}
	// check project list length
	if len(r.Projects) == 0 {
		// create a new project based on given params
		project := r.New(startConfig)
		// Add to projects list
		r.Add(project)
		// save config
		if !startConfig.NoConfig {
			err = r.Settings.Write(r)
			if err != nil {
				return err
			}
		}
	}
	// Start web server
	if r.Server.Status {
		r.Server.Parent = r
		err = r.Server.Start()
		if err != nil {
			return err
		}
		err = r.Server.OpenURL()
		if err != nil {
			return err
		}
	}
	// start workflow
	return r.Start()
}
