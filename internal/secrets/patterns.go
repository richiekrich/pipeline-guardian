// Package secrets provides functionality for detecting accidentally committed
// secrets and credentials in source code and configuration files.
package secrets

import (
	"regexp"
)

// Patterns contains regex patterns for detecting various types of credentials.
//
// Each pattern is associated with a descriptive name and is designed to detect
// a specific type of secret or credential that might be accidentally committed
// to source control.
//
// The patterns are focused on common credential formats found in:
// - Cloud provider credentials (AWS, Azure, GCP)
// - Authentication tokens (API keys, OAuth, JWT)
// - Database connection strings and passwords
// - Private keys and certificates
// - Hard-coded secrets in configuration files
//
// To extend this list with custom patterns, you can either:
// 1. Add them directly to this map
// 2. Create a configuration file with additional patterns
var Patterns = map[string]*regexp.Regexp{
	// Cloud Provider Credentials
	"AWS Access Key":          regexp.MustCompile(`(?i)aws_access_key(?:_id)?[\s]*=[\s]*['"]?([0-9a-zA-Z/+]{20,40})['"]?`),
	"AWS Secret Key":          regexp.MustCompile(`(?i)aws_secret(?:_access_key)?[\s]*=[\s]*['"]?([0-9a-zA-Z/+]{40,80})['"]?`),
	"Azure Connection String": regexp.MustCompile(`(?i)DefaultEndpointsProtocol=https;AccountName=[^;]+;AccountKey=[^;]+;EndpointSuffix=`),
	"Google API Key":          regexp.MustCompile(`(?i)AIza[0-9A-Za-z\\-_]{35}`),

	// Version Control & CI/CD Tokens
	"GitHub Token": regexp.MustCompile(`(?i)github[_\-]?token[\s]*=[\s]*['"]?([0-9a-zA-Z_\-]{35,40})['"]?`),
	"GitLab Token": regexp.MustCompile(`(?i)gitlab[_\-]?token[\s]*=[\s]*['"]?([0-9a-zA-Z_\-]{20,64})['"]?`),
	"NPM Token":    regexp.MustCompile(`(?i)(?:NPM_TOKEN|npm_token)[\s]*=[\s]*['"]?([0-9a-zA-Z]{36})['"]?`),

	// Generic Credentials
	"Generic API Key":     regexp.MustCompile(`(?i)api[_\-]?key[\s]*=[\s]*['"]?([0-9a-zA-Z_\-]{20,80})['"]?`),
	"Private Key":         regexp.MustCompile(`-----BEGIN (?:RSA|DSA|EC|OPENSSH) PRIVATE KEY-----`),
	"Password Assignment": regexp.MustCompile(`(?i)(?:password|passwd|pwd)[\s]*=[\s]*['"]([^'"]{8,})['"]`),

	// Database Credentials
	"Connection String": regexp.MustCompile(`(?i)(?:connection_string|connectionstring)[\s]*=[\s]*['"](.+?)['"]`),
	"Database Password": regexp.MustCompile(`(?i)(?:mongodb|postgres|mysql|database).*(?:password|pwd)[\s]*=[\s]*['"]([^'"]{3,})['"]`),

	// Authentication Tokens
	"JWT Token":             regexp.MustCompile(`eyJ[a-zA-Z0-9\-_]+\.eyJ[a-zA-Z0-9\-_]+\.[a-zA-Z0-9\-_]+`),
	"SSH URL with password": regexp.MustCompile(`ssh://.*:.*@.*`),
	"Basic Auth":            regexp.MustCompile(`Authorization: Basic [a-zA-Z0-9+/=]+`),
	"Bearer Token":          regexp.MustCompile(`Authorization: Bearer [a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+\.[a-zA-Z0-9\-_=]+`),
}
