/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/hyperf/go-stringable/stringable"
	"github.com/limingxinleo/go-gen/config"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
	"strings"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen {key} {name=xxx}",
	Short: "Generate code by config.json",
	Long: `Generate code by config.json. For example:
go-gen gen dao name=UserDao 
`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("failed to get current directory: %v", err)
		}

		if len(args) < 1 {
			log.Fatal("The key is required")
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalf("failed to read falgs: %v", err)
		}

		key := args[0]

		stub := config.GetCodeConfig(key)
		code := stub.CodeStub
		params := initParams(args[1:])

		for key, value := range params {
			code = strings.ReplaceAll(code, fmt.Sprintf("{%s}", key), value)
		}

		name, ok := params["name"]
		if !ok {
			log.Fatal("The params name=xxx is required")
		}

		file := path.Join(dir, stub.Path, stringable.Snake(name)+".go")

		_ = os.MkdirAll(path.Dir(file), 0755)

		_, err = os.ReadFile(file)
		if err == nil && !force {
			log.Fatalf("%s already exists", file)
		}

		err = os.WriteFile(file, []byte(code), 0644)

		if err != nil {
			log.Fatalf("failed to create %s: %v", file, err)
		}
	},
}

func initParams(args []string) map[string]string {
	var result = make(map[string]string)
	for _, arg := range args {
		value := strings.Split(arg, "=")
		if len(value) != 2 {
			log.Fatal("The params must be xx:xx")
		}

		result[value[0]] = value[1]
	}

	return result
}

func init() {
	rootCmd.AddCommand(genCmd)
	genCmd.Flags().BoolP("force", "f", false, "Whether to overwrite existing file")
}
