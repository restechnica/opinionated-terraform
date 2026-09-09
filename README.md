# Opinionated Terraform - OTF

A lightweight opinionated wrapper around Terraform.

## Intro

Over the years we found ourselves doing same thing all over again for Terraform projects at scale. It focuses heavily 
on Terraform's native partial backends, and simply makes it less cumbersome to work with. The goal is to make it as easy
to use as Terraform's native workspace feature, but without the [workspace limitations](#terraform-workspace-limitations).

## Usage

```
otf <env> <command> [terraform args...]
```

### Flags

`otf` has only two flags of its own:

```
-v, --verbose    Set log level to info
-d, --debug      Set log level to debug
```

Everything else is passed to Terraform.

### Examples

```bash
# Plan for production
otf prod plan

# Apply to staging with a specific target
otf staging apply -target module.db

# Destroy a specific resource in dev
otf dev destroy -target aws_instance.example

# Apply with auto-approve
otf prod apply -auto-approve

# Run state commands (no var-file injected)
otf prod state list
otf prod output

# Explicit init (useful after changing backend config files)
otf prod init

# All terraform flags work as expected
otf prod plan -out tfplan
otf prod apply tfplan
```

## Terraform workspace limitations

Terraform workspaces were designed for lightweight environment isolation within a single backend, but they fall short
in multi-account setups:

- **Shared credentials** - all workspaces use the same provider and backend credentials, so switching workspaces does
  not force a credential change. Nothing stops you from applying staging changes with production credentials.
- **Single backend** - workspaces share one backend configuration. You cannot point `dev` at one S3 bucket in account A
  and `prod` at another bucket in account B without re-initializing.
- **State co-location** - all workspace states live in the same backend, separated only by a key prefix. A
  misconfiguration or permission issue can expose one environment's state to another.
- **No credential boundary** - because there is no backend swap, there is no natural checkpoint that forces you to
  confirm you have the right AWS profile, role, or session active.

The partial backend approach solves all of these by treating each environment as a completely separate backend
initialization, which inherently requires the correct credentials for each target account.

## What `otf` tries to solve

The partial backend workflow is repetitive and error-prone:

```bash
terraform init -backend-config ./backends/prod.tf -reconfigure
terraform plan -var-file ./variables/prod.tfvars
terraform apply -var-file ./variables/prod.tfvars
```

Every environment switch means re-typing the right backend, the right var-file, and making sure you didn't mix them up.
`otf` automates exactly this. It re-initializes when the environment changes, injects the right `-var-file`, and passes
everything else straight through to Terraform.

It is a convenience layer on top of Terraform's own
[partial backend configuration](https://developer.hashicorp.com/terraform/language/backend#partial-configuration),
which HashiCorp already recommends for managing multiple environments.

## Expected project layout

`otf` expects your Terraform project to follow this directory convention:

```
your-terraform-project/
  main.tf
  variables.tf
  outputs.tf
  backends/
    dev.tf
    staging.tf
    prod.tf
  variables/
    dev.tfvars
    staging.tfvars
    prod.tfvars
```

Each backend file contains partial backend configuration (bucket, key, region, etc.):

```hcl
# backends/prod.tf
bucket         = "mycompany-terraform-state-prod"
key            = "prod/terraform.tfstate"
region         = "eu-west-1"
dynamodb_table = "terraform-locks-prod"
encrypt        = true
```

Each variables file contains environment-specific variable values:

```hcl
# variables/prod.tfvars
environment = "prod"
instance_type = "m5.large"
min_capacity  = 3
```

### Multi-region deployments

The `<env>` argument is just a name - it does not have to be a single word. This makes it easy to manage
cross-region failover infrastructure where the same environment is deployed to multiple regions:

```
your-terraform-project/
  backends/
    prod-eu-west-1.tf
    prod-us-east-1.tf
    staging-eu-west-1.tf
  variables/
    prod-eu-west-1.tfvars
    prod-us-east-1.tfvars
    staging-eu-west-1.tfvars
```

```bash
# Deploy to primary region
otf prod-eu-west-1 apply

# Deploy failover to US
otf prod-us-east-1 apply

# Plan both and compare
otf prod-eu-west-1 plan
otf prod-us-east-1 plan
```

Each region gets its own backend (separate state, separate credentials if needed) and its own variables
(region-specific instance types, replica counts, endpoints). The backend swap ensures you never accidentally
apply EU changes to the US region.

## Installation

### Homebrew

```bash
brew tap restechnica/tap
brew install otf
```

### Go

```bash
go install github.com/restechnica/opinionated-terraform/cmd/otf@latest
```

### Binary

Download the latest binary for your platform from
[releases](https://github.com/restechnica/opinionated-terraform/releases) and place it on your `PATH`:

```bash
curl -Lo otf https://github.com/restechnica/opinionated-terraform/releases/latest/download/otf-darwin-arm64
chmod +x otf
sudo mv otf /usr/local/bin/
```


## How it works

When you run `otf prod plan`, here is what happens:

1. **Validates** that `backends/prod.tf` and `variables/prod.tfvars` exist
2. **Checks** if the environment changed since the last run (tracked in `.terraform/.otf`)
3. **Runs `terraform init`** with `-backend-config ./backends/prod.tf -reconfigure` if the environment changed
4. **Runs `terraform plan`** with `-var-file ./variables/prod.tfvars` prepended to your arguments

For commands that don't accept `-var-file` (like `state`, `output`, `fmt`, `validate`), step 4 skips the injection and passes your arguments through directly.

## What `otf` does not do

- **No module management** - use Terraform's own module system
- **No dependency ordering** - if you need cross-stack dependencies, look at Terragrunt
- **No state manipulation** - Terraform's state commands pass through as-is
- **No config generation** - write your own `.tf` files, `otf` just passes them along
- **No lock file management** - Terraform handles its own locking
- **No workspace management** - the whole point is that you don't need workspaces

`otf` is a ~200 line wrapper. If it ever feels like it is getting in the way, you can always drop back to raw `terraform` commands. That is by design.

## License

[MIT](LICENSE)
