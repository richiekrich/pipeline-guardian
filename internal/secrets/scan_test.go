package secrets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestScanDir(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pipeline-guardian-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files with various "secrets"
	testFiles := map[string][]byte{
		"aws_creds.txt": []byte(`
# AWS credentials
aws_access_key_id = "AKIAIOSFODNN7EXAMPLE"
aws_secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
`),
		"db_config.yaml": []byte(`
database:
  username: admin
  password: "super_secret_password_123"
  host: localhost
`),
		"deploy.sh": []byte(`
#!/bin/bash
# Deploy script
curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c" https://api.example.com/deploy
`),
		"safe.txt": []byte(`
# This file doesn't contain any secrets
version: 1.0.0
environment: production
debug: false
`),
		"binary.bin": []byte{0, 1, 2, 3, 4, 5}, // Mock binary file with null bytes
	}

	// Write test files
	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, content, 0644); err != nil {
			t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	// Run the scan
	findings, err := ScanDir(tempDir, []string{"*.bin"}) // Ignore binary files
	if err != nil {
		t.Fatalf("ScanDir failed: %v", err)
	}

	// Expected findings count
	expectedCount := 3 // AWS access key, AWS secret key, and JWT token
	if len(findings) != expectedCount {
		t.Errorf("Expected %d findings, got %d", expectedCount, len(findings))
	}

	// Check if specific patterns were found
	foundPatterns := make(map[string]bool)
	for _, finding := range findings {
		foundPatterns[finding.Rule] = true

		// Verify file paths are correct
		if !filepath.IsAbs(finding.File) {
			t.Errorf("Expected absolute file path, got: %s", finding.File)
		}

		// Verify line numbers are positive
		if finding.LineNum <= 0 {
			t.Errorf("Expected positive line number, got: %d", finding.LineNum)
		}

		// Verify line text is not empty
		if finding.LineText == "" {
			t.Errorf("Expected non-empty line text for finding in %s", finding.File)
		}
	}

	// Check specific patterns
	expectedPatterns := []string{"AWS Access Key", "AWS Secret Key"}
	for _, pattern := range expectedPatterns {
		if !foundPatterns[pattern] {
			t.Errorf("Expected to find pattern '%s', but didn't", pattern)
		}
	}

	// Test filtering
	filteredFindings := FilterFindings(findings, "aws_creds.txt", nil)
	if len(filteredFindings) != 2 { // Only AWS keys
		t.Errorf("Expected 2 filtered findings by filename, got %d", len(filteredFindings))
	}

	// Test filtering by rule
	ruleFilteredFindings := FilterFindings(findings, "", []string{"AWS Access Key"})
	if len(ruleFilteredFindings) != 1 {
		t.Errorf("Expected 1 filtered finding by rule, got %d", len(ruleFilteredFindings))
	}
}

func TestSanitizeLineText(t *testing.T) {
	// Create a string longer than 100 characters
	longText := "password = "
	for i := 0; i < 200; i++ {
		longText += "x"
	}

	sanitized := sanitizeLineText(longText)
	if len(sanitized) <= 100 || !strings.HasSuffix(sanitized, "...") {
		t.Errorf("Expected long text to be truncated with '...'")
	}

	shortText := "api_key = abc123"
	if sanitizeLineText(shortText) != shortText {
		t.Errorf("Short text should not be modified")
	}
}
