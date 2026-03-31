locals {
  # Per-environment repos avoid name clashes with existing account repos and separate dev vs prod.
  ecr_env = replace(var.environment, " ", "-")
}

resource "aws_ecr_repository" "frontend" {
  name                 = "go-fullstack-frontend-${local.ecr_env}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "bff" {
  name                 = "go-fullstack-bff-${local.ecr_env}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}

resource "aws_ecr_repository" "go_api" {
  name                 = "go-fullstack-go-api-${local.ecr_env}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }
}
