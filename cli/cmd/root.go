package cmd

import (
	"fmt"
	"os"
	"runtime"

	"masscode-parser/internal/utils"

	"github.com/ondrovic/common/utils/formatters"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version string
	rootCmd *cobra.Command
)

func init() {
	rootCmd = &cobra.Command{
		Use:     "mass-code-parser <path-to-db.json> [flags]",
		Short:   "Parse and export MassCode database",
		Long:    `A tool to parse and export MassCode database in various formats.`,
		Args:    cobra.ExactArgs(1),
		RunE:    run,
		Version: version,
	}

	rootCmd.Flags().BoolP("output", "o", false, "Export results to a file")
	rootCmd.Flags().StringP("output-path", "p", "", "Path for the output file (without extension)")
	rootCmd.Flags().StringP("output-type", "t", "html", "Output type: text, html, or json")

	viper.BindPFlags(rootCmd.Flags())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}

func run(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	if !utils.FileExists(filePath) {
		return fmt.Errorf("file '%s' does not exist", filePath)
	}

	outputFlag := viper.GetBool("output")
	outputPath := viper.GetString("output-path")
	outputType := viper.GetString("output-type")

	if outputFlag {
		if outputPath == "" {
			outputPath = "mass_code_export"
		}
		outputPath = utils.EnsureCorrectExtension(outputPath, outputType)
		// correct outputPath formatting
		outputPath = formatters.FormatPath(outputPath, runtime.GOOS)

		fmt.Printf("Output will be saved to: %s\n", outputPath)
	}

	db, err := utils.ParseDatabase(filePath)
	if err != nil {
		return err
	}

	fmt.Println("Parsed MassCode database successfully")
	fmt.Printf("Found %d folders and %d snippets\n\n", len(db.Folders), len(db.Snippets))

	folderMap := utils.BuildFolderMap(db.Folders)
	outputData := utils.ProcessSnippets(db.Snippets, folderMap, db.Tags)

	var output string
	switch outputType {
	case "text":
		output, err = utils.GenerateTextOutput(outputData)
	case "html":
		output, err = utils.GenerateHTMLOutput(outputData)
	case "json":
		output, err = utils.GenerateJSONOutput(outputData)
	default:
		return fmt.Errorf("unsupported output type: %s", outputType)
	}

	if err != nil {
		return fmt.Errorf("error generating output: %v", err)
	}

	if outputFlag {
		err := os.WriteFile(outputPath, []byte(output), 0644)
		if err != nil {
			return fmt.Errorf("error writing to output file: %v", err)
		}
		fmt.Printf("Results exported to %s\n", outputPath)
	} else {
		fmt.Println(output)
	}

	return nil
}
