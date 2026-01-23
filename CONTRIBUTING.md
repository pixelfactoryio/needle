# Contributing to Needle

Thank you for considering contributing to Needle! This document outlines the process and guidelines for contributing.

## Code of Conduct

Be respectful and constructive in all interactions with the project and community members.

## How to Contribute

### Reporting Issues

When reporting issues, please include:

- A clear, descriptive title
- Steps to reproduce the problem
- Expected vs actual behavior
- Your environment (OS, Go version, Needle version)
- Relevant logs or error messages

### Submitting Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following the coding standards
3. **Add tests** for any new functionality
4. **Ensure all tests pass** with `make test`
5. **Run linters** with `make lint` and fix any issues
6. **Format your code** with `make fmt`
7. **Commit your changes** following the commit message format
8. **Push to your fork** and submit a pull request

## Development Setup

### Prerequisites

- Go 1.24 or later
- Make
- golangci-lint 2.8.0 or later
- step-cli (for certificate testing)

### Local Development

1. Clone the repository:

```bash
git clone https://github.com/pixelfactoryio/needle.git
cd needle
```

2. Install dependencies:

```bash
go mod download
```

3. Run tests:

```bash
make test
```

4. Build the binary:

```bash
make build
```

5. Run locally:

```bash
./bin/needle --http-port 8080 --https-port 8443
```

### Available Make Targets

- `make build` - Build the binary
- `make test` - Run tests with coverage
- `make lint` - Run golangci-lint
- `make fmt` - Format code with gofmt
- `make vet` - Run go vet
- `make mocks` - Generate test mocks
- `make ca-test` - Create test CA certificates
- `make cert-test` - Create test certificates

## Coding Standards

### Go Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for code formatting
- Write clear, self-documenting code
- Add comments for exported functions and types
- Keep functions focused and concise

### Code Organization

- Place business logic in `internal/app/`
- Place infrastructure code in `internal/infra/`
- Write tests alongside the code they test
- Use interfaces for dependencies

### Testing

- Write table-driven tests where appropriate
- Aim for high test coverage on business logic
- Use mocks for external dependencies
- Test error paths and edge cases

### Documentation

- Document all exported functions and types
- Keep README.md up to date
- Add inline comments for complex logic
- Update CHANGELOG.md for user-facing changes

## Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/) for clear and structured commit history.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `chore`: Maintenance tasks
- `ci`: CI/CD changes
- `refactor`: Code refactoring
- `test`: Test additions or changes
- `perf`: Performance improvements

### Examples

```
feat(dns): add support for custom CoreDNS configuration

Added ability to specify custom Corefile for CoreDNS server,
allowing users to configure advanced DNS routing rules.

Closes #24
```

```
fix(tls): resolve certificate caching race condition

Fixed race condition in certificate cache that could cause
duplicate certificate generation under high load.
```

```
docs: update installation instructions

Added Docker Compose example and clarified prerequisite versions.
```

### Scope

Optional but recommended. Examples: `dns`, `tls`, `http`, `pki`, `docs`, `ci`

### Breaking Changes

For breaking changes, add `BREAKING CHANGE:` in the footer:

```
feat(config)!: change default ports to 8080/8443

BREAKING CHANGE: Default HTTP port changed from 80 to 8080,
default HTTPS port changed from 443 to 8443 to avoid requiring
root privileges.
```

## Pull Request Guidelines

### Before Submitting

- Rebase your branch on latest `main`
- Ensure all tests pass
- Ensure linters pass with no errors
- Update documentation if needed
- Add entries to CHANGELOG.md for user-facing changes

### PR Description

Include in your PR description:

- What changes were made and why
- Link to related issues
- Screenshots/examples if applicable
- Testing notes
- Breaking changes (if any)

### Review Process

- Maintainers will review your PR
- Address feedback in new commits
- Once approved, maintainer will merge
- PRs may be squashed on merge

## Getting Help

- Open an issue for bugs or feature requests
- Reach out to maintainers for questions
- Check existing issues and PRs first

## License

By contributing to Needle, you agree that your contributions will be licensed under the MIT License.
