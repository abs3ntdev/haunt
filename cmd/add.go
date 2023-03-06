package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/abs3ntdev/haunt/src/config"
	"github.com/abs3ntdev/haunt/src/haunt"
	"github.com/spf13/cobra"
)

var addConfig config.Flags

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a project by name",
	Long:  "Adds a project by name, if path is provided it will use 'cmd/name', all flags provided will be saved in the config file. By default go install and go run will be ran",
	Args:  cobra.MatchAll(cobra.ExactArgs(1)),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) >= 1 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return getPotentialProjets(toComplete), cobra.ShellCompDirectiveNoFileComp
	},
	RunE: add,
}

func getPotentialProjets(in string) []string {
	h := haunt.NewHaunt()
	err := h.Settings.Read(&h)
	if err != nil {
		return []string{}
	}
	out := []string{}
	cmdDir, err := os.ReadDir("cmd")
	if err != nil {
		return []string{}
	}
	for _, dir := range cmdDir {
		exists := false
		for _, proj := range h.Projects {
			if dir.Name() == proj.Name {
				exists = true
				continue
			}
		}
		if !exists {
			if in == "" {
				out = append(out, dir.Name())
			} else {
				if strings.HasPrefix(dir.Name(), in) {
					out = append(out, dir.Name())
				}
			}
		}
	}
	return out
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&addConfig.Path, "path", "p", "./", "Project base path")
	addCmd.Flags().BoolVarP(&addConfig.Format, "fmt", "f", false, "Enable go fmt")
	addCmd.Flags().BoolVarP(&addConfig.Vet, "vet", "v", false, "Enable go vet")
	addCmd.Flags().BoolVarP(&addConfig.Test, "test", "t", false, "Enable go test")
	addCmd.Flags().BoolVarP(&addConfig.Generate, "generate", "g", false, "Enable go generate")
	addCmd.Flags().BoolVarP(&addConfig.Install, "install", "i", true, "Enable go install")
	addCmd.Flags().BoolVarP(&addConfig.Build, "build", "b", false, "Enable go build")
	addCmd.Flags().BoolVarP(&addConfig.Run, "run", "r", true, "Enable go run")
}

// Add a project to an existing config or create a new one
func add(cmd *cobra.Command, args []string) (err error) {
	addConfig.Name = args[0]
	h := haunt.NewHaunt()
	// read a config if exist
	err = h.Settings.Read(&h)
	if err != nil {
		return err
	}
	projects := len(h.Projects)
	// create and add a new project
	h.Add(h.New(addConfig))
	if len(h.Projects) > projects {
		// update config
		err = h.Settings.Write(h)
		if err != nil {
			return err
		}
		log.Println(h.Prefix(haunt.Green.Bold("project successfully added")))
	} else {
		log.Println(h.Prefix(haunt.Green.Bold("project can't be added")))
	}
	return nil
}
