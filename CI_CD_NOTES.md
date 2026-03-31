## CI/CD pipeline (current implementation)

The project now includes an implemented GitHub Actions workflow at `.github/workflows/aws-deploy.yml`.
The current application version in the repo is `0.1.0` for both `frontend/` and `bff-node/`.

### Branch and environment mapping

| Branch | GitHub Environment | Terraform variables file | Purpose |
|--------|---------------------|--------------------------|---------|
| `dev`  | `development`       | `infra/terraform/environments/dev.tfvars`  | Shared development deployment |
| `main` | `production`        | `infra/terraform/environments/prod.tfvars` | Production deployment |

Feature work should land in `dev` first. Promotion from `dev` to `main` is the release path to production.

### Trigger behavior

- A push to `dev` runs CI and then deploys the development environment.
- A push to `main` runs CI and then deploys the production environment.
- `workflow_dispatch` is enabled for manual execution.

### What the current workflow does

1. **Continuous integration**
   - Checks out the repository.
   - Sets up Go using `go-domain/go.mod` and runs:
     - `go mod download`
     - `go vet ./...`
     - `go test ./... -count=1`
   - Sets up Node.js `20`.
   - In `frontend/`, runs:
     - `npm ci`
     - `npm run lint`
     - `npm run build`
   - In `bff-node/`, runs:
     - `npm ci`

2. **Continuous deployment**
   - Deploys only for `dev` and `main`.
   - Uses GitHub OIDC with AWS via `aws-actions/configure-aws-credentials`.
   - Sets up Terraform `1.9.0`.
   - Runs `terraform init` and `terraform apply` inside `infra/terraform`.
   - Selects the correct tfvars file based on branch:
     - `dev` -> `environments/dev.tfvars`
     - `main` -> `environments/prod.tfvars`

### Current deployment scope

- The workflow currently deploys **Terraform-managed AWS infrastructure**.
- The Terraform configuration presently validates AWS access and environment tagging through `aws_caller_identity`.
- Application artifact deployment for the React frontend, Node BFF, and Go service is **not yet implemented** in this workflow.

### Required GitHub environment configuration

Create these GitHub Environments in the repository settings:

- `development`
- `production`

For each environment, configure:

- `AWS_REGION`
- `AWS_ROLE_TO_ASSUME`

### Runtime and tool versions currently referenced

- App version: `0.1.0`
- Node.js in CI: `20`
- Go version: `1.25.6`
- Terraform version: `1.9.0`

### Summary

The repository has moved from a planned CI/CD design to an implemented branch-based AWS deployment pipeline. At this stage, CI validates the Go and frontend codebase, while CD applies environment-specific Terraform infrastructure for `development` and `production`.