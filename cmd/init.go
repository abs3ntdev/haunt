package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/abs3ntdev/haunt/src/haunt"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generates a haunt config file using sane defaults",
	Long:  "Generates a haunt config file using sane defaults, haunt will look for a main.go file and any directories inside the relative path 'cmd' and add them all as projects",
	RunE:  defaultConfig,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func defaultConfig(cmd *cobra.Command, args []string) error {
	r := haunt.NewHaunt()
	write := true
	if _, err := os.Stat(haunt.RFile); err == nil {
		fmt.Print(r.Prefix("Config file exists. Overwire? " + haunt.Magenta.Bold("[y/n] ") + haunt.Green.Bold("(n) ")))
		var overwrite string
		fmt.Scanf("%s", &overwrite)
		write = false
		switch strings.ToLower(overwrite) {
		case "y", "ye", "yes":
			write = true
		}
	}
	if write {
		r.SetDefaults()
		err := r.Settings.Write(r)
		if err != nil {
			return err
		}
		log.Println(r.Prefix(
			"Config file has successfully been saved at .haunt.yaml",
		))
		log.Println(r.Prefix(
			"Run haunt add --help to see how to add more projects",
		))
		return nil
	}
	return nil
}
