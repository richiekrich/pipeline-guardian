# Pipeline Guardian - Usage Examples

This document provides practical examples of how to use Pipeline Guardian in various scenarios.

## Basic Scanning

### Scan the current directory

```bash
pipeline-guardian scan
```

### Scan a specific directory

```bash
pipeline-guardian scan --path ./my-project
```

### Scan with detailed output

```bash
pipeline-guardian scan --path ./my-project --output text
```

## Advanced Scanning

### Scan only GitHub Actions workflows

```bash
pipeline-guardian scan --path ./.github/workflows --type github-actions
```

### Scan with custom ignore patterns

```bash
pipeline-guardian scan --ignore "*.log,*.tmp,node_modules,test/fixtures"
```

### Output results in JSON format

```bash
pipeline-guardian scan --output json > security-report.json
```

### Output results in CSV format for import into Excel

```bash
pipeline-guardian scan --output csv > security-report.csv
```

## Integration Examples

### Pre-commit hook

Add this to your `.git/hooks/pre-commit` file:

```bash
#!/bin/bash
pipeline-guardian scan --path .
if [ $? -ne 0 ]; then
  echo "Error: Potential secrets found. Please review and fix before committing."
  exit 1
fi
```

### GitHub Actions Workflow

```yaml
name: Security Scan

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.18'
          
      - name: Install Pipeline Guardian
        run: go install github.com/richiekrich/pipeline-guardian@latest
        
      - name: Scan for secrets
        run: pipeline-guardian scan --output json > security-report.json
        
      - name: Upload security report
        uses: actions/upload-artifact@v3
        with:
          name: security-report
          path: security-report.json
```

### GitLab CI Pipeline

```yaml
stages:
  - security

security-scan:
  stage: security
  image: golang:1.18
  script:
    - go install github.com/richiekrich/pipeline-guardian@latest
    - pipeline-guardian scan --output json > security-report.json
  artifacts:
    paths:
      - security-report.json
    expire_in: 1 week
```

### Jenkins Pipeline

```groovy
pipeline {
    agent {
        docker {
            image 'golang:1.18'
        }
    }
    stages {
        stage('Security Scan') {
            steps {
                sh 'go install github.com/richiekrich/pipeline-guardian@latest'
                sh 'pipeline-guardian scan --output json > security-report.json'
            }
            post {
                always {
                    archiveArtifacts artifacts: 'security-report.json', fingerprint: true
                }
            }
        }
    }
}
```
