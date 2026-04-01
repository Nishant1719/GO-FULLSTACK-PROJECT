# Central naming for this stack (AWS + Terraform best practice: one place, env-aware).
# Pattern: go-fullstack-<environment>-<component>
#
# env_slug: safe for AWS names (no spaces). Keep in sync with var.environment.
# stack_id: project + environment prefix for Name tags and human-readable IDs.

locals {
  env_slug = replace(trimspace(var.environment), " ", "-")
  stack_id = "go-fullstack-${local.env_slug}"
}
