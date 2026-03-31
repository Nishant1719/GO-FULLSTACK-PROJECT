#!/usr/bin/env python3
"""Send AWS SSM Run Command to EC2 to deploy docker-compose.aws.yml with ECR images and RDS URL from SSM."""
from __future__ import annotations

import base64
import json
import os
import shlex
import subprocess
import sys
import time


def main() -> None:
    compose_path = os.environ["COMPOSE_FILE"]
    instance_id = os.environ["EC2_INSTANCE_ID"]
    region = os.environ["AWS_REGION"]
    ecr_registry = os.environ["ECR_REGISTRY"]
    param_name = os.environ["DATABASE_URL_PARAMETER_NAME"]
    go_api = os.environ["GO_API_IMAGE"]
    bff = os.environ["BFF_IMAGE"]
    fe = os.environ["FRONTEND_IMAGE"]

    compose_b64 = base64.b64encode(open(compose_path, "rb").read()).decode("ascii")

    commands = [
        "set -eux",
        "cd /opt/go-fullstack",
        f"export AWS_REGION={shlex.quote(region)}",
        f"export GO_API_IMAGE={shlex.quote(go_api)}",
        f"export BFF_IMAGE={shlex.quote(bff)}",
        f"export FRONTEND_IMAGE={shlex.quote(fe)}",
        "export DATABASE_URL=$(aws ssm get-parameter --name "
        + shlex.quote(param_name)
        + " --with-decryption --query Parameter.Value --output text --region "
        + shlex.quote(region)
        + ")",
        "aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin "
        + shlex.quote(ecr_registry),
        f"echo {compose_b64} | base64 -d > docker-compose.aws.yml",
        'echo "GO_API_IMAGE=$GO_API_IMAGE" > .env',
        'echo "BFF_IMAGE=$BFF_IMAGE" >> .env',
        'echo "FRONTEND_IMAGE=$FRONTEND_IMAGE" >> .env',
        'echo "DATABASE_URL=$DATABASE_URL" >> .env',
        "docker compose -f docker-compose.aws.yml pull",
        "docker compose -f docker-compose.aws.yml up -d",
    ]

    params = {"commands": commands}
    proc = subprocess.run(
        [
            "aws",
            "ssm",
            "send-command",
            "--instance-ids",
            instance_id,
            "--document-name",
            "AWS-RunShellScript",
            "--parameters",
            json.dumps(params),
            "--comment",
            "go-fullstack docker compose deploy",
            "--output",
            "json",
        ],
        capture_output=True,
        text=True,
        check=False,
    )
    if proc.returncode != 0:
        sys.stderr.write(proc.stderr or proc.stdout or "")
        raise SystemExit(proc.returncode)

    out = json.loads(proc.stdout)
    cmd_id = out["Command"]["CommandId"]
    print(f"CommandId={cmd_id}")

    deadline = time.time() + 900
    status = ""
    while time.time() < deadline:
        inv = subprocess.run(
            [
                "aws",
                "ssm",
                "get-command-invocation",
                "--command-id",
                cmd_id,
                "--instance-id",
                instance_id,
                "--output",
                "json",
            ],
            capture_output=True,
            text=True,
            check=True,
        )
        data = json.loads(inv.stdout)
        status = data.get("Status", "")
        if status in ("Success", "Cancelled", "Failed", "TimedOut"):
            print(f"Status={status}")
            if status != "Success":
                print(data.get("StandardErrorContent", "") or data.get("StandardOutputContent", ""))
                raise SystemExit(1)
            return
        time.sleep(5)

    raise SystemExit("Timed out waiting for SSM command")


if __name__ == "__main__":
    main()
