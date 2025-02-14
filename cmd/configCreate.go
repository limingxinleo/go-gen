package cmd

import (
	"fmt"
	"github.com/limingxinleo/go-gen/config"
	"github.com/spf13/cobra"
)

// configCreateCmd represents the configCreate command
var configCreateCmd = &cobra.Command{
	Use:   "config:create",
	Short: "Create config in current project",
	Long:  `Create config in current project`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetCodeConfig("dao")
		fmt.Println(conf)
	},
}

func init() {
	rootCmd.AddCommand(configCreateCmd)
}
