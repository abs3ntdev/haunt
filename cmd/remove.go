package cmd

import (
	"log"
	"strings"

	"github.com/abs3ntdev/haunt/src/haunt"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [names]",
	Short: "Removes all projects by name from config file",
	Args:  cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getProjectNames(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: remove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func getProjectNames(input string) []string {
	r := haunt.NewHaunt()
	// read a config if exist
	err := r.Settings.Read(&r)
	if err != nil {
		return []string{}
	}
	names := []string{}
	for _, project := range r.Schema.Projects {
		if strings.HasPrefix(project.Name, input) {
			names = append(names, project.Name)
		}
	}
	return names
}

// Remove a project from an existing config
func remove(cmd *cobra.Command, args []string) (err error) {
	r := haunt.NewHaunt()
	// read a config if exist
	err = r.Settings.Read(&r)
	if err != nil {
		return err
	}
	for _, arg := range args {
		err := r.Schema.Remove(arg)
		if err != nil {
			log.Println(r.Prefix(haunt.Red.Bold(arg + " project not found")))
			continue
		}
		log.Println(r.Prefix(haunt.Green.Bold(arg + " successfully removed")))
	}
	// update config
	err = r.Settings.Write(r)
	if err != nil {
		return err
	}

	return nil
}
