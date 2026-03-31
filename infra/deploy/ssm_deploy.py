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


def _aws_json(args: list[str]) -> dict:
    proc = subprocess.run(args, capture_output=True, text=True, check=False)
    if proc.returncode != 0:
        sys.stderr.write(proc.stderr or proc.stdout or "")
        raise SystemExit(proc.returncode)
    return json.loads(proc.stdout)


def wait_for_instance_running(instance_id: str, region: str, timeout_sec: int = 300) -> None:
    """Fail fast if the instance is missing/terminated; wait until state is running."""
    deadline = time.time() + timeout_sec
    while time.time() < deadline:
        data = _aws_json(
            [
                "aws",
                "ec2",
                "describe-instances",
                "--instance-ids",
                instance_id,
                "--region",
                region,
                "--output",
                "json",
            ]
        )
        reservations = data.get("Reservations") or []
        if not reservations or not reservations[0].get("Instances"):
            raise SystemExit(
                f"No EC2 instance data for {instance_id} in {region} (wrong id or region?)"
            )
        inst = reservations[0]["Instances"][0]
        state = inst["State"]["Name"]
        if state == "running":
            print(f"EC2 {instance_id} state=running", flush=True)
            return
        if state in ("terminated", "shutting-down"):
            raise SystemExit(
                f"EC2 {instance_id} is {state}; cannot deploy. Recreate the instance (e.g. terraform apply)."
            )
        print(f"EC2 {instance_id} state={state} (waiting...)", flush=True)
        time.sleep(10)
    raise SystemExit(f"Timeout: EC2 {instance_id} did not reach running within {timeout_sec}s")


def wait_for_ssm_online(instance_id: str, region: str, timeout_sec: int = 600) -> None:
    """SendCommand returns InvalidInstanceId until SSM lists the instance with PingStatus=Online."""
    deadline = time.time() + timeout_sec
    while time.time() < deadline:
        data = _aws_json(
            [
                "aws",
                "ssm",
                "describe-instance-information",
                "--filters",
                f"Key=InstanceIds,Values={instance_id}",
                "--region",
                region,
                "--output",
                "json",
            ]
        )
        items = data.get("InstanceInformationList") or []
        if items:
            ping = items[0].get("PingStatus", "")
            if ping == "Online":
                print(f"SSM PingStatus=Online for {instance_id}", flush=True)
                return
            print(f"SSM PingStatus={ping} for {instance_id} (waiting...)", flush=True)
        else:
            print(
                f"SSM has no InstanceInformation yet for {instance_id} (agent registering...)",
                flush=True,
            )
        time.sleep(15)
    raise SystemExit(
        f"Timeout after {timeout_sec}s: SSM never reported Online for {instance_id}. "
        "Check: IAM instance profile includes AmazonSSMManagedInstanceCore, VPC endpoints or "
        "public egress for SSM, and instance is running."
    )


def main() -> None:
    compose_path = os.environ["COMPOSE_FILE"]
    instance_id = os.environ["EC2_INSTANCE_ID"]
    region = os.environ["AWS_REGION"]
    ecr_registry = os.environ["ECR_REGISTRY"]
    param_name = os.environ["DATABASE_URL_PARAMETER_NAME"]
    go_api = os.environ["GO_API_IMAGE"]
    bff = os.environ["BFF_IMAGE"]
    fe = os.environ["FRONTEND_IMAGE"]

    wait_for_instance_running(instance_id, region)
    wait_for_ssm_online(instance_id, region)

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
