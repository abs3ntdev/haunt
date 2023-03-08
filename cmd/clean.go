package cmd

import (
	"log"

	"github.com/abs3ntdev/haunt/src/haunt"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:     "clean",
	Aliases: []string{"c"},
	Short:   "Deletes the haunt config file",
	RunE:    clean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

// Clean remove haunt file
func clean(cmd *cobra.Command, args []string) (err error) {
	h := haunt.NewHaunt()
	if err := h.Settings.Remove(haunt.HFile); err != nil {
		return err
	}
	log.Println(h.Prefix(haunt.Green.Bold("config file removed successfully removed")))
	return nil
}
