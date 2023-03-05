package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "haunt",
	Short: `Haunt: Go task automater and live-reloader.
 .-.
(o o)
| O \
 \   \
  '~~'
	`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		logrus.Error(err)
		return
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
