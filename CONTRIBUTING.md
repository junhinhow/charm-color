# Contributing to charm-color

Thank you for your interest in contributing!

## How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/new-function`)
3. Add your function to `color.go` with proper doc comments
4. Write tests in `color_test.go`
5. Run tests: `go test ./...`
6. Commit using conventional commits (`feat:`, `fix:`, `docs:`)
7. Open a Pull Request

## Code Style

- All public functions must accept and return `color.Color` interface
- Clamp numeric inputs to valid ranges (don't panic)
- Run `gofmt` before committing
- All exported symbols must have doc comments

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
