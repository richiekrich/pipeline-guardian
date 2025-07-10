/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/richiekrich/pipeline-guardian/internal/secrets"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan pipeline configuration files for security issues",
	Long: `The scan command analyzes your pipeline configuration files 
(GitHub Actions, GitLab CI, Jenkins, etc.) for security vulnerabilities 
and compliance issues, including accidentally committed secrets or credentials.

Examples:
  pipeline-guardian scan --path ./github/workflows
  pipeline-guardian scan --type github-actions --output json`,
	Run: func(cmd *cobra.Command, args []string) {
		scanPath, _ := cmd.Flags().GetString("path")
		scanType, _ := cmd.Flags().GetString("type")
		outputFormat, _ := cmd.Flags().GetString("output")
		ignorePatterns, _ := cmd.Flags().GetStringSlice("ignore")

		// Verify path exists
		if _, err := os.Stat(scanPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Path '%s' does not exist\n", scanPath)
			os.Exit(1)
		}

		fmt.Println("📊 Scanning for security issues...")
		fmt.Printf("Path: %s\n", scanPath)
		fmt.Printf("Type: %s\n", scanType)
		fmt.Printf("Output format: %s\n", outputFormat)

		// Perform the secret/credential scan
		findings, err := secrets.ScanDir(scanPath, ignorePatterns)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning for secrets: %v\n", err)
			os.Exit(1)
		}

		// Filter findings based on scan type if specified
		if scanType != "auto" && scanType != "" {
			// Apply scan type-specific filtering (e.g., only GitHub Actions files)
			filteredFindings := []secrets.Finding{}
			for _, f := range findings {
				// For GitHub Actions, focus on workflow files
				if scanType == "github-actions" && (strings.HasSuffix(f.File, ".yml") || strings.HasSuffix(f.File, ".yaml")) {
					filteredFindings = append(filteredFindings, f)
				} else if scanType == "gitlab-ci" && strings.HasSuffix(f.File, ".gitlab-ci.yml") {
					filteredFindings = append(filteredFindings, f)
				} else if scanType == "jenkins" && strings.HasSuffix(f.File, "Jenkinsfile") {
					filteredFindings = append(filteredFindings, f)
				} else if scanType == "all" {
					filteredFindings = append(filteredFindings, f)
				}
			}
			findings = filteredFindings
		}

		// Output findings based on format
		outputFindings(findings, outputFormat)
	},
}

// outputFindings prints the findings in the specified format
func outputFindings(findings []secrets.Finding, format string) {
	if len(findings) == 0 {
		fmt.Println("✅ Scan completed. No security issues found.")
		return
	}

	fmt.Printf("🔴 Found %d potential security issues:\n", len(findings))

	switch format {
	case "json":
		// Output as JSON
		jsonOutput, _ := json.MarshalIndent(findings, "", "  ")
		fmt.Println(string(jsonOutput))

	case "csv":
		// Output as CSV
		fmt.Println("File,Rule,LineNumber,LineContent")
		for _, f := range findings {
			// Escape quotes in line content for CSV
			lineContent := strings.ReplaceAll(f.LineText, "\"", "\"\"")
			fmt.Printf("\"%s\",\"%s\",%d,\"%s\"\n", f.File, f.Rule, f.LineNum, lineContent)
		}

	default: // "text" format
		// Output as formatted text
		for i, f := range findings {
			fmt.Printf("%d) %s (line %d)\n", i+1, f.Rule, f.LineNum)
			fmt.Printf("   File: %s\n", f.File)
			fmt.Printf("   Content: %s\n\n", f.LineText)
		}
	}

	fmt.Println("\n⚠️ Warning: Review these potential issues and ensure no sensitive information is committed.")
	fmt.Println("   Remember to add any false positives to your .gitignore or use a tool like git-secrets.")
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Define flags for the scan command
	scanCmd.Flags().StringP("path", "p", ".", "Path to the directory containing pipeline configuration files")
	scanCmd.Flags().StringP("type", "t", "auto", "Type of pipeline (github-actions, gitlab-ci, jenkins, all, etc.)")
	scanCmd.Flags().StringP("output", "o", "text", "Output format (text, json, csv)")
	scanCmd.Flags().StringSliceP("ignore", "i", []string{".git", "node_modules", "vendor", "*.jpg", "*.png", "*.gif"}, "Patterns to ignore")
}
