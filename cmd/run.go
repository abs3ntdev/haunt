package cmd

import (
	"log"
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

// haunt workflow
func run(cmd *cobra.Command, args []string) (err error) {
	r := haunt.NewHaunt()

	// read a config if exist
	err = r.Settings.Read(&r)
	if err != nil {
		return err
	}
	if len(args) >= 1 {
		// filter by name flag if exist
		r.Projects = r.Filter(args)
		if len(r.Projects) == 0 {
			log.Println(r.Prefix("No valid project found, exiting. Check your config file or run haunt add"))
			return
		}
	}
	// increase file limit
	if r.Settings.FileLimit != 0 {
		if err = r.Settings.Flimit(); err != nil {
			return err
		}
	}

	//  web server
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

	// run workflow
	return r.Run()
}
