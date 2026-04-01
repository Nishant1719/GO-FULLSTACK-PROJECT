resource "aws_ecr_repository" "frontend" {
  name                 = "go-fullstack-frontend-${local.env_slug}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${local.stack_id}-ecr-frontend"
  }
}

resource "aws_ecr_repository" "bff" {
  name                 = "go-fullstack-bff-${local.env_slug}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${local.stack_id}-ecr-bff"
  }
}

resource "aws_ecr_repository" "go_api" {
  name                 = "go-fullstack-go-api-${local.env_slug}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = "${local.stack_id}-ecr-go-api"
  }
}
