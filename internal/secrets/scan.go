package secrets

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

// Finding represents a detected credential leak in a file
type Finding struct {
	File     string // Path to the file containing the leak
	Rule     string // Name of the pattern rule that matched
	LineNum  int    // Line number where the leak was found
	LineText string // Content of the line (potentially truncated for display)
	Offset   []int  // Start and end position of the match in the file content
}

// ScanDir recursively scans a directory for credential leaks
func ScanDir(rootPath string, ignorePatterns []string) ([]Finding, error) {
	var findings []Finding

	// Walk through all files in the directory
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip files with errors
		}

		// Skip directories
		if info.IsDir() {
			// Check if directory should be ignored
			for _, pattern := range ignorePatterns {
				if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check if file should be ignored
		for _, pattern := range ignorePatterns {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				return nil
			}
		}

		// Skip binary files and large files
		if info.Size() > 1024*1024*5 { // 5MB limit
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// Skip if likely binary file
		if isBinary(content) {
			return nil
		}

		// Map to keep track of lines where findings were already reported
		reportedLines := make(map[int]bool)

		// Check content against each pattern
		for ruleName, pattern := range Patterns {
			matches := pattern.FindAllIndex(content, -1)
			for _, loc := range matches {
				// Get line number for this match
				lineNum, lineText := getLineInfo(content, loc[0])

				// Skip if we already reported an issue on this line
				if reportedLines[lineNum] {
					continue
				}
				reportedLines[lineNum] = true

				// Create a finding
				findings = append(findings, Finding{
					File:     path,
					Rule:     ruleName,
					LineNum:  lineNum,
					LineText: sanitizeLineText(lineText),
					Offset:   loc,
				})
			}
		}

		return nil
	})

	return findings, err
}

// isBinary does a simple check if a file appears to be binary
func isBinary(content []byte) bool {
	// Check first 512 bytes for NULL bytes which indicates a binary file
	end := 512
	if len(content) < end {
		end = len(content)
	}

	for _, b := range content[:end] {
		if b == 0 {
			return true
		}
	}
	return false
}

// getLineInfo returns the line number and line text for a position in a file
func getLineInfo(content []byte, position int) (int, string) {
	// Count newlines before the position to get the line number (1-based)
	lineNum := 1
	for i := 0; i < position; i++ {
		if content[i] == '\n' {
			lineNum++
		}
	}

	// Extract the line text
	scanner := bufio.NewScanner(bytes.NewReader(content))
	currentLine := 0
	for scanner.Scan() {
		currentLine++
		if currentLine == lineNum {
			return lineNum, scanner.Text()
		}
	}

	return lineNum, ""
}

// sanitizeLineText removes sensitive data and truncates if needed
func sanitizeLineText(text string) string {
	// Truncate long lines
	maxLen := 100
	if len(text) > maxLen {
		return text[:maxLen] + "..."
	}

	return text
}

// FilterFindings filters findings based on criteria
func FilterFindings(findings []Finding, filePattern string, ruleNames []string) []Finding {
	if filePattern == "" && len(ruleNames) == 0 {
		return findings
	}

	filtered := []Finding{}
	for _, f := range findings {
		// Filter by file pattern
		includeFile := true
		if filePattern != "" {
			matched, _ := filepath.Match(filePattern, filepath.Base(f.File))
			includeFile = matched
		}

		// Filter by rule names
		includeRule := true
		if len(ruleNames) > 0 {
			includeRule = false
			for _, rule := range ruleNames {
				if strings.EqualFold(f.Rule, rule) {
					includeRule = true
					break
				}
			}
		}

		if includeFile && includeRule {
			filtered = append(filtered, f)
		}
	}

	return filtered
}
