package cmd

import (
	"log"

	"github.com/abs3ntdev/haunt/src/haunt"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Deletes the haunt config file",
	RunE:  clean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

// Clean remove haunt file
func clean(cmd *cobra.Command, args []string) (err error) {
	r := haunt.NewHaunt()
	if err := r.Settings.Remove(haunt.RFile); err != nil {
		return err
	}
	log.Println(r.Prefix(haunt.Green.Bold("config file removed successfully removed")))
	return nil
}
