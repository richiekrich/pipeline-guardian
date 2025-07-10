# Pipeline Guardian

<div align="center">
  <img src="https://img.shields.io/badge/version-1.0.0-blue.svg" alt="Version 1.0.0">
  <img src="https://img.shields.io/badge/license-MIT-green.svg" alt="License MIT">
  <img src="https://img.shields.io/badge/go-%3E%3D%201.18-blue.svg" alt="Go Version">
</div>

<div align="center">
  <h3>🛡️ Your DevSecOps sidekick 🚀</h3>
  <p>A powerful CLI tool for securing your CI/CD pipelines</p>
</div>

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [Basic Commands](#basic-commands)
  - [Scanning for Secrets](#scanning-for-secrets)
  - [Configuration](#configuration)
  - [Output Formats](#output-formats)
  - [Filtering Results](#filtering-results)
  - [Ignoring Files](#ignoring-files)
- [Examples](#examples)
- [Best Practices](#best-practices)
- [Contribution](#contribution)
- [License](#license)

## Overview

Pipeline Guardian is a security-focused CLI tool designed to help developers and DevOps teams implement security best practices in their CI/CD pipelines. It provides tools for scanning, validating, and generating reports on security issues in your pipeline configurations, with a special focus on preventing accidental exposure of credentials and secrets.

## Features

- **Secret and credential detection**: Find accidentally committed API keys, passwords, tokens and other sensitive information
- **Pipeline validation**: Validate your pipeline configurations against security best practices
- **Support for multiple CI/CD platforms**: Works with GitHub Actions, GitLab CI, Jenkins, and more
- **Flexible output formats**: Get results in human-readable text, JSON, or CSV formats
- **Customizable scanning**: Include/exclude specific files or directories from scanning
- **Easy integration**: Designed to work both in local environments and as part of automated CI/CD processes

## Installation

### Prerequisites

- Go 1.18 or higher

### From source

```bash
git clone https://github.com/richiekrich/pipeline-guardian.git
cd pipeline-guardian
go install
```

### Using Go

```bash
go install github.com/richiekrich/pipeline-guardian@latest
```

## Usage

### Basic Commands

Pipeline Guardian offers several commands:

```
pipeline-guardian [command] [flags]
```

Available commands:

- `scan`: Scan pipeline configuration files for security issues
- `validate`: Validate pipeline configurations against security policies
- `report`: Generate security reports for your CI/CD pipelines

For help on any command:

```
pipeline-guardian [command] --help
```

### Scanning for Secrets

The `scan` command searches for accidental exposure of credentials and secrets in your codebase or pipeline configuration files.

**Basic usage**:

```
pipeline-guardian scan --path ./my-project
```

This will recursively scan all files in the specified directory for potential secrets and credentials.

**Command options**:

```
Flags:
  -h, --help             help for scan
  -i, --ignore strings   Patterns to ignore (default [.git,node_modules,vendor,*.jpg,*.png,*.gif])
  -o, --output string    Output format (text, json, csv) (default "text")
  -p, --path string      Path to the directory containing pipeline configuration files (default ".")
  -t, --type string      Type of pipeline (github-actions, gitlab-ci, jenkins, all, etc.) (default "auto")
```

### Configuration

Pipeline Guardian can use a configuration file to customize its behavior. By default, it looks for a file named `.pipeline-guardian.yaml` in your home directory.

To specify a custom config file:

```
pipeline-guardian --config /path/to/config.yaml scan
```

Example configuration file:

```yaml
# .pipeline-guardian.yaml
ignore:
  - "*.backup"
  - "*.min.js"
  - "vendor/**"

rules:
  # Enable/disable specific secret detection rules
  aws-keys: true
  github-tokens: true
  private-keys: true

severity:
  # Configure minimum severity level
  level: "medium" # Options: low, medium, high, critical
```

### Output Formats

Pipeline Guardian supports multiple output formats to fit your workflow:

#### Text Output (default)

Human-readable format that's easy to scan visually:

```
pipeline-guardian scan --output text
```

Example output:
```
🔴 Found 2 potential security issues:

1) AWS Access Key (line 22)
   File: config/deploy.yml
   Content: aws_access_key_id = "AKIAIOSFODNN7EXAMPLE"

2) Private Key (line 45)
   File: scripts/deploy.sh
   Content: -----BEGIN RSA PRIVATE KEY-----
```

#### JSON Output

Machine-readable format for integration with other tools:

```
pipeline-guardian scan --output json
```

Example output:
```json
[
  {
    "File": "config/deploy.yml",
    "Rule": "AWS Access Key",
    "LineNum": 22,
    "LineText": "aws_access_key_id = \"AKIAIOSFODNN7EXAMPLE\"",
    "Offset": [462, 504]
  },
  {
    "File": "scripts/deploy.sh",
    "Rule": "Private Key",
    "LineNum": 45,
    "LineText": "-----BEGIN RSA PRIVATE KEY-----",
    "Offset": [1255, 1285]
  }
]
```

#### CSV Output

For importing into spreadsheets or other data tools:

```
pipeline-guardian scan --output csv
```

Example output:
```
File,Rule,LineNumber,LineContent
"config/deploy.yml","AWS Access Key",22,"aws_access_key_id = \"AKIAIOSFODNN7EXAMPLE\""
"scripts/deploy.sh","Private Key",45,"-----BEGIN RSA PRIVATE KEY-----"
```

### Filtering Results

You can focus your scan on specific types of pipeline configurations:

```
pipeline-guardian scan --type github-actions
```

This will only report issues found in GitHub Actions workflow files (`.yml` or `.yaml` files).

Available pipeline types:
- `github-actions`: GitHub Actions workflow files
- `gitlab-ci`: GitLab CI configuration files
- `jenkins`: Jenkins pipeline files
- `all`: All file types
- `auto`: Auto-detect (default)

### Ignoring Files

To exclude specific files or directories from being scanned:

```
pipeline-guardian scan --ignore "*.log,backup/*,node_modules"
```

By default, Pipeline Guardian ignores the following patterns:
- `.git`
- `node_modules`
- `vendor`
- `*.jpg`
- `*.png`
- `*.gif`

## Examples

### Basic secret scanning

```bash
# Scan current directory
pipeline-guardian scan

# Scan specific directory
pipeline-guardian scan --path ./my-project

# Scan with JSON output
pipeline-guardian scan --path ./my-project --output json
```

### CI/CD platform-specific scanning

```bash
# Scan GitHub Actions workflows
pipeline-guardian scan --type github-actions --path ./.github/workflows

# Scan GitLab CI configuration
pipeline-guardian scan --type gitlab-ci --path ./

# Scan Jenkins pipeline files
pipeline-guardian scan --type jenkins --path ./jenkins
```

### Advanced usage

```bash
# Scan with custom ignore patterns
pipeline-guardian scan --ignore "*.backup,*.lock,node_modules,dist"

# Validate pipeline against security policies
pipeline-guardian validate --policy strict --pipeline ./.github/workflows/deploy.yml

# Generate a security report
pipeline-guardian report --format pdf --output ./reports
```

## Best Practices

1. **Regular scans**: Run Pipeline Guardian as part of your pre-commit hooks to prevent secrets from being committed.
2. **CI integration**: Include Pipeline Guardian in your CI/CD process to catch secrets that might have been missed.
3. **Custom patterns**: Extend the default patterns with your organization-specific patterns to catch custom secrets.
4. **Use alongside other tools**: Combine Pipeline Guardian with other security tools like SAST scanners for comprehensive security coverage.
5. **Act on findings**: Always review and address the issues found by Pipeline Guardian. Never ignore security warnings!

## Detected Secret Types

Pipeline Guardian can detect various types of secrets and credentials:

| Secret Type | Description | Example Pattern |
|-------------|-------------|----------------|
| AWS Access Keys | AWS access key IDs | `aws_access_key_id = "AKIA..."` |
| AWS Secret Keys | AWS secret access keys | `aws_secret_access_key = "wJal..."` |
| GitHub Tokens | GitHub personal access tokens | `github_token = "ghp_..."` |
| Generic API Keys | Various API keys | `api_key = "..."` |
| Private Keys | RSA, DSA, EC, or OpenSSH private keys | `-----BEGIN RSA PRIVATE KEY-----` |
| Passwords | Password assignments in code | `password = "secret123"` |
| Connection Strings | Database connection strings | `connection_string = "..."` |
| Database Passwords | Database-specific password assignments | `postgres_password = "..."` |
| JWT Tokens | JSON Web Tokens | `eyJhbGciOiJ...` |
| SSH URLs with passwords | SSH URLs that include passwords | `ssh://user:pass@host` |
| Basic Auth Headers | HTTP Basic Authentication headers | `Authorization: Basic ZGV2...` |
| Bearer Tokens | HTTP Bearer Authentication tokens | `Authorization: Bearer eyJ...` |

## Contribution

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

<div align="center">
  <p>Built with ❤️ by <a href="https://github.com/richiekrich">Richie K Rich</a></p>
</div>
