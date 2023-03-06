package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abs3ntdev/haunt/src/haunt"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run haunt, optionally provide the name of projects to only run those otherwise will run all configured projects",
	Args:  cobra.MatchAll(cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getProjectNamesToRun(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: run,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func getProjectNamesToRun(input string) []string {
	h := haunt.NewHaunt()
	// read a config if exist
	err := h.Settings.Read(&h)
	if err != nil {
		return []string{}
	}
	names := []string{}
	for _, project := range h.Projects {
		if strings.HasPrefix(project.Name, input) {
			names = append(names, project.Name)
		}
	}
	return names
}

// haunt workflow
func run(cmd *cobra.Command, args []string) (err error) {
	h := haunt.NewHaunt()

	// read a config if exist
	err = h.Settings.Read(&h)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println(h.Prefix("No config file found, initializing one for you"))
			err = defaultConfig(cmd, args)
			if err != nil {
				log.Println(h.Prefix("Failed to generate default config: " + err.Error()))
			}
			err = h.Settings.Read(&h)
			if err != nil {
				return fmt.Errorf(h.Prefix("Failed to read config file: " + err.Error()))
			}
		} else {
			return err
		}
	}
	if len(args) >= 1 {
		// filter by name flag if exist
		h.Projects = h.Filter(args)
		if len(h.Projects) == 0 {
			log.Println(h.Prefix("No valid project found, exiting. Check your config file or run haunt add"))
			return
		}
	}
	// increase file limit
	if h.Settings.FileLimit != 0 {
		if err = h.Settings.Flimit(); err != nil {
			return err
		}
	}

	//  web server
	if h.Server.Status {
		h.Server.Parent = h
		err = h.Server.Start()
		if err != nil {
			return err
		}
		err = h.Server.OpenURL()
		if err != nil {
			return err
		}
	}

	// run workflow
	return h.Run()
}
