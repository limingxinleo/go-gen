package cmd

import (
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/limingxinleo/go-gen/src/json2struct"
	"github.com/marhaupe/json2struct/pkg/generator"
	"github.com/marhaupe/json2struct/pkg/parse"
	"github.com/spf13/cobra"
)

type Json2StructFlags struct {
	String             string
	File               string
	ShouldUseClipboard bool
}

var (
	json2StructFlags Json2StructFlags

	json2StructCmd = &cobra.Command{
		Use:     "json2struct",
		Short:   "json2struct generates Go type definitions for a JSON",
		Version: version,
		Args:    cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			Run(json2StructFlags)
		},
	}
)

func init() {
	json2StructCmd.Flags().StringVarP(&json2StructFlags.String, "string", "s", "", "JSON string")
	json2StructCmd.Flags().StringVarP(&json2StructFlags.File, "file", "f", "", "path to JSON file")
	json2StructCmd.Flags().BoolVarP(&json2StructFlags.ShouldUseClipboard, "clipboard", "c", false, "read from and write types to clipboard")

	rootCmd.AddCommand(json2StructCmd)
}

func Run(flags Json2StructFlags) {
	var userInputNode parse.Node
	var err error
	reader := json2struct.NewReader()

	switch {
	case flags.ShouldUseClipboard:
		var userInput string
		userInput, err = clipboard.ReadAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
		userInputNode, err = reader.ParseFromString(userInput)
	case flags.File != "":
		userInput := reader.ReadFromFile(flags.File)
		userInputNode, err = reader.ParseFromString(userInput)
	case flags.String != "":
		userInputNode, err = reader.ParseFromString(flags.String)
	default:
		userInputNode, err = reader.ReadFromEditor()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if userInputNode == nil {
		return
	}

	output, err := generator.GenerateOutputFromAST(userInputNode)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	fmt.Println(output)

	if flags.ShouldUseClipboard {
		err = clipboard.WriteAll(output)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
		fmt.Println("\nSaved output to clipboard")
	}
}
