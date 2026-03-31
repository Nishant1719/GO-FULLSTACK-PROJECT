variable "environment" {
  description = "Logical environment (development or production)."
  type        = string
}

variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "vpc_cidr" {
  description = "CIDR for the demo VPC."
  type        = string
  default     = "10.0.0.0/16"
}

variable "ec2_instance_type" {
  description = "EC2 instance type for Docker Compose host."
  type        = string
  default     = "t3.micro"
}

variable "rds_instance_class" {
  description = "RDS instance class for PostgreSQL."
  type        = string
  default     = "db.t3.micro"
}

variable "rds_allocated_storage" {
  description = "RDS allocated storage in GB."
  type        = number
  default     = 20
}

variable "rds_skip_final_snapshot" {
  description = "Skip final snapshot on RDS destroy (set false for production data retention)."
  type        = bool
  default     = true
}

variable "rds_backup_retention_days" {
  description = "RDS backup retention in days (0 disables backups)."
  type        = number
  default     = 1
}

variable "db_name" {
  description = "PostgreSQL database name."
  type        = string
  default     = "go_domain_db"
}

variable "db_username" {
  description = "PostgreSQL master username."
  type        = string
  default     = "postgres"
}
