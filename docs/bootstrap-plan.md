# Plan: `otf bootstrap <env>` command

## Context

When a Terraform project manages its own state backend (e.g., S3 bucket + DynamoDB table), there's a chicken-and-egg problem: you need the bucket before remote state works, but the bucket is Terraform-managed. The manual workaround is: init with local state, apply, migrate state, clean up. This is painful in CI/CD pipelines and tedious locally.

`otf bootstrap <env>` automates the entire flow as a single command. It's provider-agnostic - the user writes their own Terraform config, declares an `otf_backend` output, and `otf` orchestrates everything: creating resources, generating the backend config files, migrating state, and optionally pushing the generated files to git.

## User's Terraform project setup

The user's project needs:
- A `variables/<env>.tfvars` file (as usual)
- NO `backend` block in the terraform config (bootstrap creates it)
- An `otf_backend` output:

```hcl
output "otf_backend" {
  value = {
    type           = "s3"
    bucket         = aws_s3_bucket.state.id
    region         = "eu-west-1"
    dynamodb_table = aws_dynamodb_table.lock.id
    key            = "terraform.tfstate"
  }
}
```

## The bootstrap flow

```
otf bootstrap prod [--push]
```

Steps executed in sequence (fail-fast):

1. `terraform init -input=false` - plain init, local state (no backend block exists)
2. `terraform apply -var-file variables/<env>.tfvars -auto-approve -input=false` - create backend resources
3. `terraform output -json otf_backend` - read backend config from outputs
4. Generate `backends/<env>.tf` - partial backend config (key-value pairs without `type`)
5. Generate `backend.tf` - the `terraform { backend "<type>" {} }` block
6. `terraform init -backend-config backends/<env>.tf -migrate-state -input=false` - migrate local state to remote
7. Cleanup - remove `terraform.tfstate` and `terraform.tfstate.backup`, write `.terraform/.otf`
8. (if `--push`) `git add backends/<env>.tf backend.tf`, `git commit`, `git push`

## The `--push` flag

Three modes based on usage:

- **No flag (default)**: generate files, log what was created. Developer commits when ready (local workstation) or pipeline acts on the log output.
- **`--push`**: commit the generated files and push. For CI/CD where changes need to flow back to the repo. Commit message: `chore: bootstrap backend for <env>`.

Git operations shell out to `git` via the existing Commander, following the same pattern as semverbot's `pkg/git/cli.go`. No Go git library - relies on the environment's git credentials (GitHub Actions checkout, SSH keys, etc.).

## Files to create

### `pkg/backend/generate.go`
Parse the `otf_backend` output map and generate two files:
- `ParseOutput(raw map[string]interface{}) (Config, error)` - extracts `type`, validates, converts remaining values to strings
- `WritePartialConfig(env string, cfg Config) error` - writes `backends/<env>.tf` with `key = "value"` lines, creates `backends/` dir if needed
- `WriteBackendBlock(cfg Config) error` - writes `backend.tf` with `terraform { backend "<type>" {} }`

### `pkg/backend/generate_test.go`
Tests for parsing (valid, missing type, nested values rejected) and file generation (content correctness, directory creation, overwrite idempotency). Use `t.TempDir()`.

### `pkg/git/api.go`
Interface following semverbot's pattern:
```go
type API interface {
    Add(files ...string) error
    Commit(message string) error
    Push() error
}
```

### `pkg/git/cli.go`
Implementation shelling out via Commander: `git add <files>`, `git commit -m <msg>`, `git push`.

### `pkg/git/cli_test.go`
Tests using a mock Commander to verify correct git args.

### `pkg/core/bootstrap.go`
Orchestrator function `Bootstrap(tf terraform.API, git git.API, env string, push bool) error` that runs all steps in sequence. Reuses `tf.Run()` for apply (which already injects `-var-file`).

### `pkg/core/bootstrap_test.go`
Tests with mock `terraform.API` and mock `git.API`: happy path call order, each-step-failure stops subsequent steps, push=false skips git, push=true calls git add/commit/push.

### `pkg/cli/bootstrap/bootstrap.go`
Cobra subcommand: `cobra.ExactArgs(1)`, `--push` bool flag, delegates to `core.Bootstrap()`.

## Files to modify

### `pkg/commander/commander.go`
Fix `Output` to separate stdout from stderr. Currently both go into one buffer, which breaks `terraform output -json` parsing. Change to capture stdout and stderr in separate buffers, return stdout only.

### `pkg/terraform/api.go`
Add three methods:
- `InitSimple() error` - `terraform init -input=false`
- `InitMigrateState(env string) error` - `terraform init -backend-config ... -migrate-state -input=false`
- `Output(name string) (map[string]interface{}, error)` - `terraform output -json <name>`, parse JSON envelope

### `pkg/terraform/cli.go`
Implement the three new methods. `Output` uses `Commander.Output()` (stdout-only after the fix), parses the JSON `{"value": {...}}` envelope, returns the inner map.

### `pkg/env/env.go`
Add `ValidateVariables(env string) error` - only checks `variables/<env>.tfvars` exists (bootstrap can't check `backends/<env>.tf` since it doesn't exist yet).

### `pkg/cli/defaults.go`
Add constants: `DefaultBackendFile = "backend.tf"`, `DefaultOutputName = "otf_backend"`.

### `pkg/cli/root/root.go`
Import bootstrap package, register with `cmd.AddCommand(bootstrap.NewCommand(tf))`.

## Implementation order

1. `pkg/commander/commander.go` - fix Output stdout/stderr separation
2. `pkg/cli/defaults.go` + `pkg/env/env.go` - new constants + ValidateVariables
3. `pkg/terraform/api.go` + `pkg/terraform/cli.go` - interface + implementation
4. `pkg/backend/generate.go` + tests - new package (parallel with step 3)
5. `pkg/git/api.go` + `pkg/git/cli.go` + tests - new package (parallel with step 3-4)
6. `pkg/core/bootstrap.go` + tests - orchestrator (depends on 2-5)
7. `pkg/cli/bootstrap/bootstrap.go` + `pkg/cli/root/root.go` - wire up CLI

## Verification

1. `go build ./...` - compiles
2. `go test ./...` - all tests pass
3. `go vet ./...` - no issues
4. Manual test with a real Terraform project that creates an S3 bucket to confirm the full flow