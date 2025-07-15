package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "wxgw",
	Short: "wxgw",
	Long:  `微信企业网关`,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
