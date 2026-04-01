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


def wait_for_ssm_online(instance_id: str, region: str, timeout_sec: int = 1200) -> None:
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
        "Check: IAM instance profile has AmazonSSMManagedInstanceCore; instance has outbound "
        "internet (public IP or NAT) or VPC endpoints for SSM; on the instance: "
        "`sudo systemctl status amazon-ssm-agent` and `/var/log/amazon/ssm/amazon-ssm-agent.log`."
    )


def resolve_ecr_and_images(region: str) -> tuple[str, str, str, str]:
    """Fill missing registry / image URIs (e.g. re-run deploy-only loses upstream job outputs)."""
    ecr_registry = os.environ.get("ECR_REGISTRY", "").strip()
    if not ecr_registry:
        proc = subprocess.run(
            [
                "aws",
                "sts",
                "get-caller-identity",
                "--query",
                "Account",
                "--output",
                "text",
            ],
            capture_output=True,
            text=True,
            check=True,
        )
        account = proc.stdout.strip()
        ecr_registry = f"{account}.dkr.ecr.{region}.amazonaws.com"
        print(f"ECR_REGISTRY was empty; using {ecr_registry}", flush=True)

    go_api = os.environ.get("GO_API_IMAGE", "").strip()
    bff = os.environ.get("BFF_IMAGE", "").strip()
    fe = os.environ.get("FRONTEND_IMAGE", "").strip()
    if not go_api or not bff or not fe:
        image_tag = os.environ.get("IMAGE_TAG", "").strip()
        deploy_env = os.environ.get("DEPLOY_ENV", "").strip()
        if not image_tag or not deploy_env:
            raise SystemExit(
                "ECR image env vars are empty (common when re-running only the deploy job). "
                "Use 'Re-run all jobs', or set GO_API_IMAGE, BFF_IMAGE, FRONTEND_IMAGE, "
                "and ECR_REGISTRY — or set DEPLOY_ENV + IMAGE_TAG for automatic URIs."
            )
        go_api = f"{ecr_registry}/go-fullstack-go-api-{deploy_env}:{image_tag}"
        bff = f"{ecr_registry}/go-fullstack-bff-{deploy_env}:{image_tag}"
        fe = f"{ecr_registry}/go-fullstack-frontend-{deploy_env}:{image_tag}"
        print(f"Image URIs were empty; using tag {image_tag} for env {deploy_env}", flush=True)

    return ecr_registry, go_api, bff, fe


def main() -> None:
    compose_path = os.environ["COMPOSE_FILE"]
    instance_id = os.environ["EC2_INSTANCE_ID"]
    region = os.environ["AWS_REGION"]
    param_name = os.environ["DATABASE_URL_PARAMETER_NAME"]

    ecr_registry, go_api, bff, fe = resolve_ecr_and_images(region)

    wait_for_instance_running(instance_id, region)
    wait_for_ssm_online(instance_id, region)

    compose_b64 = base64.b64encode(open(compose_path, "rb").read()).decode("ascii")

    # SSM can report Online before cloud-init finishes: /opt may be missing and Docker may not exist yet.
    commands = [
        "set -eux",
        "mkdir -p /opt/go-fullstack",
        "chown ec2-user:ec2-user /opt/go-fullstack",
        "cd /opt/go-fullstack",
        # Up to ~15m: user_data may still be installing Docker when SSM is already Online.
        "i=0; while [ $i -lt 90 ]; do command -v docker >/dev/null 2>&1 && docker info >/dev/null 2>&1 && break; echo Waiting for Docker...; sleep 10; i=$((i+1)); done",
        "command -v docker",
        "docker info",
        # AL2023 repos may not ship docker-compose-plugin; ensure `docker compose` exists (dnf or GitHub binary).
        "dnf install -y docker-compose-plugin || true",
        (
            "if ! docker compose version >/dev/null 2>&1; then "
            "mkdir -p /usr/local/lib/docker/cli-plugins && "
            "ARCH=$(uname -m) && "
            'case "$ARCH" in x86_64) DC_ARCH=x86_64 ;; aarch64) DC_ARCH=aarch64 ;; '
            '*) echo "unsupported arch: $ARCH"; exit 1 ;; esac && '
            "curl -fsSL "
            '"https://github.com/docker/compose/releases/download/v2.29.7/docker-compose-linux-${DC_ARCH}" '
            "-o /usr/local/lib/docker/cli-plugins/docker-compose && "
            "chmod +x /usr/local/lib/docker/cli-plugins/docker-compose; "
            "fi"
        ),
        "docker compose version",
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
