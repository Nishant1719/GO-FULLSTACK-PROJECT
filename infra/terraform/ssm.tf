resource "aws_iam_role_policy" "ec2_read_db_url" {
  name_prefix = "go-fullstack-dburl-"
  role        = aws_iam_role.ec2.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Sid    = "ReadDatabaseUrlParameter"
      Effect = "Allow"
      Action = [
        "ssm:GetParameter",
        "ssm:GetParameters"
      ]
      Resource = aws_ssm_parameter.database_url.arn
    }]
  })
}

resource "aws_ssm_parameter" "database_url" {
  name        = "/go-fullstack/${var.environment}/database_url"
  description = "PostgreSQL connection string for Go API (demo)"
  type        = "SecureString"
  # Adopt an existing parameter (e.g. left over from a prior partial apply) instead of failing with ParameterAlreadyExists.
  overwrite = true
  value = format(
    "postgresql://%s:%s@%s:%s/%s?sslmode=require",
    var.db_username,
    urlencode(random_password.db.result),
    aws_db_instance.main.address,
    aws_db_instance.main.port,
    var.db_name
  )

  depends_on = [aws_db_instance.main]
}
