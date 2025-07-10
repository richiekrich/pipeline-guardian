/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pipeline-guardian",
	Short: "A DevSecOps tool for securing your CI/CD pipelines",
	Long: `Pipeline Guardian is your DevSecOps sidekick, designed to help you 
implement security best practices in your CI/CD pipelines.

It provides tools for scanning, monitoring, and enforcing security 
policies throughout your development and deployment processes.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🛡️ Pipeline Guardian: your DevSecOps sidekick 🚀")
		fmt.Println("")
		fmt.Println("Quick Start:")
		fmt.Println("  pipeline-guardian scan            - Scan for secrets in current directory")
		fmt.Println("  pipeline-guardian scan --path ./  - Scan specific directory")
		fmt.Println("  pipeline-guardian validate        - Validate pipeline configurations")
		fmt.Println("  pipeline-guardian report          - Generate security reports")
		fmt.Println("")
		fmt.Println("Get Help:")
		fmt.Println("  pipeline-guardian --help          - Show all available commands")
		fmt.Println("  pipeline-guardian scan --help     - Show options for the scan command")
		fmt.Println("")
		fmt.Println("For detailed documentation, visit: https://github.com/richiekrich/pipeline-guardian")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pipeline-guardian.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".pipeline-guardian" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".pipeline-guardian")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
