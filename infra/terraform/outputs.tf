output "account_id" {
  value = data.aws_caller_identity.current.account_id
}

output "environment" {
  value = var.environment
}

output "stack_id" {
  description = "Stack prefix (project + environment) used in names and tags."
  value       = local.stack_id
}

output "ec2_instance_id" {
  value = aws_instance.app.id
}

output "ec2_public_ip" {
  value = aws_instance.app.public_ip
}

output "database_url_parameter_name" {
  value = aws_ssm_parameter.database_url.name
}

output "ecr_registry_url" {
  value = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com"
}

output "app_url" {
  value = "http://${aws_instance.app.public_ip}"
}
