package cmd

import (
	"log"

	"github.com/abs3ntdev/haunt/src/haunt"
	"github.com/spf13/cobra"
)

var Version = "v0.2.11"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints current verison",
	Run:   version,
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Version print current version
func version(cmd *cobra.Command, args []string) {
	h := haunt.NewHaunt()
	log.Println(h.Prefix(haunt.Green.Bold(Version)))
}
