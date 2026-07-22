# Repository Guidelines

## Project Structure & Module Organization

`autojournal/` is the Wails desktop application. `internal/` is organised by responsibility:

- `domain/` contains entities and repository interfaces.
- `storage/` contains SQLite/sqlx implementations.
- `service/` contains validation and application workflows.
- `handler/` exposes services to the Wails frontend.
- `migrations/` contains ordered SQL migrations such as `001_init.up.sql`.

The Vue/TypeScript frontend is in `autojournal/frontend/`: code in `src/`, assets in `src/assets/`, bindings in `wailsjs/` (do not edit).

## Build, Test, and Development Commands

Run commands from `autojournal/`:

```bash
make dev    # start the Wails app with frontend development tooling
make build  # build the desktop application
make test   # run all Go tests
make fmt    # format Go sources
make vet    # run static checks
make check  # format, vet, and test before submitting work
```

For frontend-only work: `cd frontend && npm run dev`; use `npm run build` to type-check and build Vue.

## Coding Style & Naming Conventions

Use `gofmt`; do not hand-format Go code. Packages are lowercase; exported identifiers use `PascalCase`, unexported ones `camelCase`. Keep database access in `storage`, business rules in `service`, and Wails code in `handler`.

Use descriptive migration names with a zero-padded sequence, for example `003_add_reminders.up.sql`. Vue components use `PascalCase` filenames (for example, `VehicleForm.vue`); TypeScript variables and functions use `camelCase`.

## Testing Guidelines

Place Go tests beside covered code as `*_test.go`, using `TestThing_Condition`. Prefer service tests with fake repositories; use SQLite tests for SQL behaviour. Run `make check` before a pull request. Test new business rules and bug fixes.

## Teaching and Collaboration Mode

Act as a Go backend teacher and reviewer, not a default implementer. For Go, database, architecture, repositories, services, Wails bindings, and structure: explain the goal, name the file, give small steps, flag mistakes, and provide only a minimal unblocker example. Review the user's attempt before the next iteration.

Write backend code only when explicitly asked, a small reference is needed, after a learning attempt, or for a low-value technical correction. Prioritise understanding over speed.

Frontend work may be implemented directly. Explain changed files, UI/backend data flow, state, and component structure.

## Commit & Pull Request Guidelines

Existing history uses short, imperative Russian descriptions, often scoped to a phase (for example, `Прописал миграции`). Make one focused change per commit. PRs must explain the change, checks run, and related issue or phase step; add frontend screenshots. Never commit databases, credentials, or generated build output.
