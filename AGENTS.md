# Repository Guidelines

## Project Structure & Module Organization

File Browser combines a Go backend with a Vue 3/TypeScript frontend. `main.go` starts the application; backend concerns are split into focused packages such as `cmd/`, `http/`, `auth/`, `files/`, `storage/`, `share/`, and `users/`. Go tests live beside their implementations as `*_test.go`. Frontend code is under `frontend/src/`, with static files in `frontend/public/` and unit tests in `__tests__/`. Documentation sources and site configuration live in `www/docs/` and `www/mkdocs.yml`; branding assets are in `branding/`.

## Build, Test, and Development Commands

- `task build` installs frontend dependencies, builds the frontend, and compiles the Go binary.
- `go run .` runs the backend locally; use `go mod download` first on a fresh clone.
- `cd frontend && pnpm install && pnpm run dev` starts the Vite development server.
- `go test --race ./...` runs all backend tests with race detection, matching CI.
- `cd frontend && pnpm run test` runs Vitest once; `pnpm run lint` checks frontend code.
- `task docs` builds documentation into `www/public/`; `task docs:serve` serves it on port 8000.

Use Go 1.26, Node.js 24+, pnpm 10+, and Task where possible to match CI.

## Coding Style & Naming Conventions

Format Go files with `gofmt`; use idiomatic package names, exported `PascalCase` identifiers, and unexported `camelCase` identifiers. Keep packages cohesive and errors contextual. Frontend code follows ESLint and Prettier (`pnpm run lint`, `pnpm run format`), with two-space indentation, double quotes, and trailing commas where supported. Name Vue components in `PascalCase.vue` and TypeScript utilities in descriptive lowercase filenames.

## Testing Guidelines

Add tests close to changed behavior. Use table-driven Go tests when multiple cases share setup, naming functions `TestXxx`. Place frontend tests in a nearby `__tests__/` directory and name them `*.test.ts`. No fixed coverage threshold is enforced, but regressions and edge cases should accompany fixes. Run both backend and frontend suites before submitting cross-stack changes.

## Commit & Pull Request Guidelines

Follow Conventional Commits used in history and enforced for PR titles: `feat: ...`, `fix(scope): ...`, `docs: ...`, or `chore: ...`. Keep commits focused and subjects imperative. Target `master`; complete the PR template with a clear description, linked issue (`Fixes #123`), test evidence, and screenshots for visible UI changes. The project is maintenance-only, and translation changes must go through Transifex.
