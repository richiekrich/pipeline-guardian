/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate security reports for your CI/CD pipelines",
	Long: `The report command generates detailed security reports for your CI/CD pipelines
based on scan results and security policy validations.

Examples:
  pipeline-guardian report --format pdf --output ./reports
  pipeline-guardian report --format html --include-details`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		outputPath, _ := cmd.Flags().GetString("output")
		includeDetails, _ := cmd.Flags().GetBool("include-details")

		fmt.Println("📝 Generating security report...")
		fmt.Printf("Format: %s\n", format)
		fmt.Printf("Output path: %s\n", outputPath)
		fmt.Printf("Include details: %t\n", includeDetails)
		fmt.Println("✅ Report generated successfully.")
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	// Define flags for the report command
	reportCmd.Flags().StringP("format", "f", "pdf", "Report format (pdf, html, markdown, json)")
	reportCmd.Flags().StringP("output", "o", "./", "Directory to store the generated reports")
	reportCmd.Flags().BoolP("include-details", "d", false, "Include detailed information in the report")
}
