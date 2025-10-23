# CI/CD Pipeline Documentation

This document describes the CI/CD pipeline setup for the dep-parser project.

## Overview

The CI/CD pipeline consists of multiple workflows that ensure code quality, security, and proper deployment of the dep-parser project.

## Workflows

### 1. Quick Tests (`.github/workflows/test.yml`)

**Purpose:** Fast feedback on pull requests and pushes
**Triggers:** 
- Push to main/master branches
- Pull request events

**Jobs:**
- **Quick Unit Tests:** Runs `go test -short ./...` for fast feedback
- **Timeout:** 10 minutes

### 2. CI/CD Pipeline (`.github/workflows/ci.yml`)

**Purpose:** Comprehensive testing and deployment pipeline
**Triggers:**
- Push to main/master branches
- Pull requests to main/master branches

**Jobs:**

#### Lint and Format
- Runs `make lint` (gofmt, golangci-lint, goimports)
- Checks `go mod tidy`
- Validates license compliance
- Must pass before other jobs run

#### Test
- Runs unit tests with coverage on multiple Go versions (1.23, 1.24)
- Uploads coverage reports to Codecov
- Requires lint job to pass

#### Security Scan
- Runs Gosec security scanner
- Runs Trivy vulnerability scanner
- Uploads results to GitHub Security tab
- Continues on error (doesn't fail the build)

#### Build
- Builds the project on multiple platforms (Linux, Windows, macOS)
- Tests build artifacts
- Matrix testing across different OS and Go versions

#### Integration Tests
- Runs integration tests (if any exist)
- Only runs on push events (not PRs)
- Searches for `*_integration_test.go` files

#### Release
- Creates GitHub releases when tags are pushed
- Builds binaries for multiple platforms
- Generates release notes automatically
- Only runs on version tags (`v*.*.*`)

#### Deploy Documentation
- Deploys documentation to GitHub Pages
- Only runs on main/master branch
- Sets up Node.js for documentation tools

### 3. Dependency Updates (`.github/workflows/dependency-updates.yml`)

**Purpose:** Automated dependency management and security monitoring
**Triggers:**
- Daily schedule (6 AM UTC)
- Manual workflow dispatch

**Jobs:**

#### Update Go Dependencies
- Runs `go get -u ./...` to update dependencies
- Commits changes automatically
- Creates PRs for dependency updates
- Only pushes to main branch

#### Check Security
- Runs `govulncheck` to find known vulnerabilities
- Creates GitHub issues if vulnerabilities are found
- Labels issues with 'security', 'dependencies', 'automated'

### 4. Code Quality (`.github/workflows/quality.yml`)

**Purpose:** Code quality monitoring and performance tracking
**Triggers:**
- Push to main/master branches
- Pull requests
- Weekly schedule (Sundays at 2 AM UTC)

**Jobs:**

#### Code Quality Checks
- Runs staticcheck for static analysis
- Runs ineffassign to find ineffective assignments
- Runs misspell for spelling checks
- Fails on TODO/FIXME/XXX comments in non-test files

#### Performance Benchmarks
- Runs benchmark tests
- Stores benchmark results using benchmark-action
- Alerts on performance regression (>200% threshold)
- Comments on alerts with performance issues

#### Code Complexity Analysis
- Runs gocyclo to analyze function complexity
- Generates complexity reports
- Uploads reports as artifacts
- Flags functions with complexity > 15

#### Documentation Coverage
- Checks documentation coverage using goreportcard
- Validates godoc generation
- Reports undocumented public functions

## Environment Variables

### Global Variables
- `GO_VERSION`: '1.24' (default Go version)
- `COVERAGE_THRESHOLD`: 80 (minimum test coverage percentage)

### Required Secrets

#### For Release Workflow
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions

#### For Documentation Deployment
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions

#### For Dependency Updates
- `GITHUB_TOKEN`: Automatically provided by GitHub Actions

## Matrix Testing

The pipeline uses matrix testing for:

### Go Versions
- **Test Job:** 1.23, 1.24
- **Build Job:** 1.24

### Operating Systems
- **Build Job:** ubuntu-latest, windows-latest, macos-latest

### Build Targets (Release)
- Linux AMD64/ARM64
- Darwin AMD64/ARM64
- Windows AMD64

## Quality Gates

### Test Coverage
- Minimum coverage: 80%
- Enforced by `make unit` command

### Linting
- All linting must pass
- Includes gofmt, golangci-lint, goimports
- License compliance check

### Security
- Security scans run but don't fail the build
- Vulnerabilities create GitHub issues for tracking

### Performance
- Performance regression threshold: 200%
- Automated alerts and comments on performance issues

## Artifacts

### Generated Artifacts
- **Complexity Reports:** Code complexity analysis
- **Benchmark Results:** Performance benchmark data
- **Build Binaries:** Multi-platform binaries (on release)
- **Coverage Reports:** Test coverage data

### Coverage Reporting
- Unit test coverage uploaded to Codecov
- Coverage threshold enforced in CI

## Security Scanning

### Tools Used
- **Gosec:** Go security scanner
- **Trivy:** Vulnerability scanner for dependencies
- **govulncheck:** Go vulnerability database checker

### Security Workflow
1. Scan code and dependencies for vulnerabilities
2. Upload results to GitHub Security tab
3. Create issues for critical vulnerabilities
4. Continue build execution (security warnings don't fail build)

## Release Process

### Automated Releases
- Triggered by version tags (`v*.*.*`)
- Builds multi-platform binaries
- Creates GitHub release with notes
- Uploads binaries as release assets

### Release Artifacts
- `dep-parser-linux-amd64`
- `dep-parser-linux-arm64`
- `dep-parser-darwin-amd64`
- `dep-parser-darwin-arm64`
- `dep-parser-windows-amd64.exe`

## Monitoring and Alerts

### Performance Monitoring
- Benchmark results tracked over time
- Automated alerts on performance regression
- GitHub comments on performance issues

### Quality Monitoring
- Code complexity trends
- Test coverage tracking
- Documentation coverage monitoring

### Security Monitoring
- Daily vulnerability scanning
- Automated issue creation for security issues
- Dependency update automation

## Local Development

### Running CI Checks Locally
```bash
# Bootstrap dependencies
make bootstrap

# Run all checks
make all

# Run specific checks
make lint          # Linting and formatting
make unit          # Unit tests with coverage
make check-licenses # License compliance
```

### Testing Workflows Locally
```bash
# Install act for local GitHub Actions testing
brew install act

# Run workflows locally
act -W .github/workflows/ci.yml
act -W .github/workflows/quality.yml
```

## Troubleshooting

### Common Issues

#### Build Failures
- Check Go version compatibility
- Ensure all dependencies are available
- Verify `make bootstrap` completed successfully

#### Test Failures
- Check test coverage threshold (80%)
- Verify test files are properly named
- Check for missing test dependencies

#### Lint Failures
- Run `make lint-fix` to auto-fix formatting
- Check for unused imports or variables
- Verify license compliance

#### Security Scan Failures
- Review vulnerability reports
- Update affected dependencies
- Check for false positives

### Debugging

#### Viewing Logs
- GitHub Actions provides detailed logs for each step
- Artifacts contain additional debugging information
- Use `continue-on-error: true` for non-critical failures

#### Local Reproduction
- Use `act` to run workflows locally
- Check matrix combinations individually
- Verify environment setup matches CI

## Maintenance

### Updating Workflows
- Review GitHub Actions versions regularly
- Update Go versions in matrix testing
- Monitor for deprecated actions

### Adding New Checks
- Follow existing job patterns
- Use appropriate matrix strategies
- Consider impact on build times

### Performance Optimization
- Use caching for Go modules and dependencies
- Parallelize independent jobs
- Use matrix testing efficiently

## Contributing

When contributing to the CI/CD pipeline:

1. Test changes locally using `act`
2. Ensure all existing workflows continue to pass
3. Update this documentation for any significant changes
4. Consider impact on build times and resources
5. Follow existing patterns and conventions
