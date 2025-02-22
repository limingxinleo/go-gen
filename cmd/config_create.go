package cmd

import (
	"github.com/limingxinleo/go-gen/config"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// configCreateCmd represents the configCreate command
var configCreateCmd = &cobra.Command{
	Use:   "config:create",
	Short: "Create config in current project",
	Long:  `Create config in current project`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current directory: %v", err)
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalf("failed to read falgs: %v", err)
		}

		err = config.CreateConfigDir(dir, force)
		if err != nil {
			log.Fatalf("failed to init config: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCreateCmd)

	configCreateCmd.Flags().BoolP("force", "f", false, "Whether to overwrite existing config")
}
