resource "random_password" "db" {
  length  = 16
  special = false
}

resource "aws_security_group" "rds" {
  name_prefix = "go-fullstack-rds-"
  description = "PostgreSQL from EC2 only"
  vpc_id      = aws_vpc.main.id

  ingress {
    description     = "Postgres from app host"
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.ec2.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_db_subnet_group" "main" {
  # name_prefix avoids collision with an orphan "go-fullstack-<env>" group left in the account from a failed run.
  name_prefix = "go-fullstack-${replace(var.environment, " ", "-")}-"
  subnet_ids  = [aws_subnet.private_a.id, aws_subnet.private_b.id]
}

resource "aws_db_instance" "main" {
  identifier = "go-fullstack-${replace(var.environment, " ", "-")}"

  engine               = "postgres"
  engine_version       = "16"
  instance_class       = var.rds_instance_class
  allocated_storage    = var.rds_allocated_storage
  storage_type         = "gp3"
  db_name              = var.db_name
  username             = var.db_username
  password             = random_password.db.result
  db_subnet_group_name = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  publicly_accessible          = false
  skip_final_snapshot          = var.rds_skip_final_snapshot
  backup_retention_period      = var.rds_backup_retention_days
  deletion_protection          = false
  apply_immediately            = true
  auto_minor_version_upgrade   = true
  performance_insights_enabled = false

  lifecycle {
    prevent_destroy = false
  }
}
