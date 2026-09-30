# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

Nine AWS Lambda functions that implement UVA Library's APTrust submission workflow, plus the
Postgres migrations for the shared submission database. Each `apt-submit-*/` directory is a
**separate Go module** (`package main`, its own `go.mod`/`Makefile`). There is no top-level module
and no `go.work` (it is gitignored).

## Build / run / deploy commands

All commands run from inside a lambda directory, e.g. `cd apt-submit-state-manager`.

```sh
make build     # symlink common sources + build local cmdline binary -> bin/cmd
make linux     # symlink common sources + build lambda binary -> bin/bootstrap + bin/deployment.zip
make all       # both
make clean     # go clean + rm -rf bin
make dep       # go get -u && go mod tidy && go mod verify
make fmt       # go fmt
make vet       # go vet
```

- `make build` hardcodes `GOOS=darwin GOARCH=amd64`; `make linux` hardcodes `GOOS=linux GOARCH=amd64`.
- The lambda binary is named `bootstrap` (AL2023 `provided.*` runtime).
- Build tags select the entrypoint: `-tags cmdline` for local, `-tags lambda.norpc,lambda` for deploy.
  Every `main()` is behind one of these tags, so a plain `go build` with no tags will not link.
- **There are no tests in this repo** — no `*_test.go` files and no `test` make target. Verification
  is `make vet` plus running the `cmdline` binary against a real database/bus.

Running locally requires the function's env vars to be exported (see its `config.go`). The
convention is an untracked `env.set` file (gitignored) that you source first.

```sh
# event-bridge-driven functions (flags come from lambda-common/main-cmdline.go)
./bin/cmd --eventname bag.accepted --cid uva --sid sid-xxx --bid bag-name --extra '{"etag":"..."}'

# API-gateway-driven functions have their own main.go with function-specific flags
./bin/cmd --cid uva              # e.g. apt-submit-submission-register
```

Deployment is AWS CodeBuild, not local: `pipeline/buildspec.yml` builds all nine and uploads
`deployment.zip` per function to `s3://aptrust-submit-lambda-deployable/<timestamp>/<fn>/` and to
`latest/`, then writes the version to SSM `/lambdas/aptrust-submit-lambda-deployable/latest`.
`pipeline/deployspec.yml` points `uva-apt-submit-<fn>-staging` at `latest/`.
**Adding or removing a lambda means editing both `buildspec.yml` and `deployspec.yml`.**

## The `lambda-common` mechanism (most important thing to understand)

`lambda-common/*.go` holds shared source that is **symlinked into each lambda directory** by the
`common:` make target, because the files are `package main` and cannot be imported. Consequences:

- The symlinked filenames (`env.go`, `events.go`, `helpers.go`, `definitions.go`, `main-*.go`,
  `s3.go`) are listed in the root `.gitignore` so the copies never get committed; only
  `lambda-common/` is tracked (`!lambda-common/*.go`).
- Each lambda's `common:` target lists only the files it actually needs, so **editing a common file
  affects every lambda that symlinks it** — check the `common:` targets before changing one.
- Adding a new shared file requires three edits: create it in `lambda-common/`, add its name to
  `.gitignore`, and add the `ln -s` line to each consuming `Makefile`.
- The symlinks are not present in a fresh checkout; `make` creates them. Editing a symlinked file
  from within a lambda directory edits the shared original.

### The `process()` contract

Each lambda supplies `process()`; `lambda-common` supplies the `main()` that calls it. Exactly one
`main-*.go` variant is symlinked per lambda, which determines both the trigger and the signature:

| Common main | Signature | Lambdas |
| --- | --- | --- |
| `main-lambda-eb.go` (EventBridge) | `process(messageId, messageSrc string, rawMsg json.RawMessage) error` | event-audit, ingest-status, purge, state-manager |
| `main-lambda-sqs.go` (SQS of EB events) | same as above, once per SQS record | mailer |
| `main-lambda-apigw.go` (API Gateway) | `process(messageId, messageSrc string, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)` | bag-status, submission-initiate, submission-register, submission-status |

`main-cmdline.go` is the EventBridge-shaped local driver. API-gateway lambdas cannot use it, so
they each carry their own `main.go` under `//go:build cmdline`.

Returning an error from an EB/SQS `process()` causes the event to be retried/requeued; unrecognised
event types therefore log `WARNING: unexpected event type ..., ignoring` and return `nil`.

## Shared external packages

Two `uvalib` modules carry the cross-service contracts; most behaviour changes start by checking them:

- `github.com/uvalib/aptrust-submit-db-dao/uvaaptsdao` — Postgres DAO (`NewDao(host, port, user, password, name)`),
  plus the `SubmissionStatus*` / `BagStatus*` constants and `Err*NotFound` sentinels.
- `github.com/uvalib/aptrust-submit-bus-definitions/uvaaptsbus` — EventBridge event names
  (`EventSubmissionApprove`, `EventBagAccepted`, …), `UvaBusEvent`, `UvaWorkflowEvent`,
  `MakeBusEvent`/`MakeWorkflowEvent`.

Each `go.mod` has commented-out `replace` directives for developing against sibling checkouts
(`../../aptrust-submit-db-dao/uvaaptsdao`, `../../aptrust-submit-bus-definitions/uvaaptsbus`) —
uncomment for local work, re-comment before committing. `make dep` is how versions get bumped
(hence the recurring "Latest dependancies" commits).

## Workflow / event flow

Events flow over an EventBridge bus; the bus name and publisher identity come from
`EVENT_BUS_NAME` / `EVENT_SRC_NAME`. If `EVENT_BUS_NAME` is empty, `NewEventBus` returns
`nil, ErrConfig` **quietly** and `publishWorkflowEvent` becomes a no-op — that is deliberate, for
running locally without emitting events.

Bag/submission assets live at `<bucket-or-efs>/<clientId>/<submissionId>/<bagId>/...`.

- **submission-register** (API GW) — creates a submission row, mints `sid-<xid>`, returns the
  deposit bucket + `<clientId>/<sid>` path.
- **submission-initiate** (API GW) — lists the deposited S3 objects, checks the discovered bag
  folders match the request, publishes `EventSubmissionValidate`.
- **state-manager** (EB) — the workflow's state machine. Handles reconcile-fail, incomplete,
  approve/approved/abandoned, and bag submitted/rejected/accepted; validates the current state
  before transitioning, and publishes `EventBagInitiate` (on approval) and `EventSubmissionComplete`
  (when the last bag is accepted). Clients with an empty `ApprovalEmail` are auto-approved.
- **mailer** (SQS) — renders `templates/*.template` (embedded via `go:embed`) for approve /
  validate-fail / reconcile-fail and sends via SMTP; skips the approval mail for auto-approve clients.
- **ingest-status** (EB, scheduled) — polls APTrust for bags in `BagStatusPendingIngest`, capped by
  `MAX_REQUESTS` (shuffles and takes a subset when over), publishes `EventBagAccepted`/`EventBagRejected`,
  and warns when an ingest has been pending more than 24h.
- **bag-status** / **submission-status** (API GW) — read-only state lookups by `bid` / `sid` query param.
- **purge** (EB) — deletes cached assets for a completed/purge-commanded submission. The S3-side
  purge is currently commented out; only the EFS cache is cleared.
- **event-audit** (EB) — logs every event on the bus; no side effects.

Some event types are intentionally *not* handled here: submission validate / validate-fail /
bag-built state updates are done by the issuing services. The scheduler and approval lambdas were
retired (AWS scheduler and state-manager took over) — don't reintroduce them.

## Database migrations

`db/migrations/NNNNN_<description>.{up,down}.sql`, sequential 5-digit prefixes, currently through
`00027`. Always write both directions. Applied with `golang-migrate` via `scripts/migrate.ksh`
(needs `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`) into the `migrations` history
table, driven in CI by `pipeline/migratespec.yml`. Migrations are shipped to the deploy bucket by
`buildspec.yml` and applied as a separate pipeline stage from the function deploys, so schema
changes and code that depends on them land at different times.

## Code conventions

- Configuration lives in each lambda's `config.go` as a `Config` struct plus
  `loadConfiguration() (*Config, error)`, built from the `lambda-common/env.go` helpers
  (`envWithDefault`, `ensureSetAndNonEmpty`, `envToInt`, `envToBool`), and echoed back with
  `[CONFIG] Name = [value]` lines (secrets printed as `[REDACTED]`). It is loaded **inside
  `process()`**, per invocation, not at cold start — keep it that way so config errors surface as
  event/HTTP failures.
- Logging is `fmt.Printf` with a severity prefix: `DEBUG:`, `INFO:`, `WARNING:`, `ERROR:`. No logging
  library. Errors are logged at the point of failure *and* returned.
- API Gateway errors go through `apiGatewayProxyErrorResponse(status, err)`; not-found conditions are
  detected with `errors.Is(err, uvaaptsdao.Err*NotFound)` — the sentinels are values, not types, so
  `errors.As` with a `*error` target matches everything and must not be used (see commit 20fff57).
- Files open with a `//`-comment banner and close with a `//\n// end of file\n//` trailer. Match it.
- Commit messages are short, sentence-case, no prefix or ticket (e.g. "Retire the approval lambda").
