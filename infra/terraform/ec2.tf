resource "aws_security_group" "ec2" {
  name_prefix = "go-fullstack-ec2-"
  description = "Docker Compose host"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "HTTP for demo app"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
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

resource "aws_iam_role" "ec2" {
  name_prefix = "go-fullstack-ec2-"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
  })
}

resource "aws_iam_role_policy_attachment" "ec2_ssm" {
  role       = aws_iam_role.ec2.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
}

resource "aws_iam_role_policy_attachment" "ec2_ecr" {
  role       = aws_iam_role.ec2.name
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
}

resource "aws_iam_instance_profile" "ec2" {
  name_prefix = "go-fullstack-ec2-"
  role        = aws_iam_role.ec2.name
}

resource "aws_instance" "app" {
  ami                    = data.aws_ami.al2023.id
  instance_type          = var.ec2_instance_type
  subnet_id              = aws_subnet.public_a.id
  vpc_security_group_ids = [aws_security_group.ec2.id]
  iam_instance_profile   = aws_iam_instance_profile.ec2.name

  user_data = <<-EOT
    #!/bin/bash
    set -eux
    dnf install -y docker
    systemctl enable --now docker
    usermod -aG docker ec2-user
    dnf install -y docker-compose-plugin
    mkdir -p /opt/go-fullstack
    chown ec2-user:ec2-user /opt/go-fullstack
  EOT

  root_block_device {
    volume_size = 30
    volume_type = "gp3"
  }

  depends_on = [
    aws_internet_gateway.main,
    aws_db_instance.main,
    aws_iam_role_policy.ec2_read_db_url,
  ]
}
