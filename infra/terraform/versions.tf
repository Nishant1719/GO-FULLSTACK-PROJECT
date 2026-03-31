terraform {
  required_version = ">= 1.5.0"

  # State is configured via -backend-config in CI (see CI_CD_NOTES.md). Local dev: copy backend.hcl.example.
  backend "s3" {}

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.6"
    }
  }
}
