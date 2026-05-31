# net_utils

This repository contains Terraform automation for provisioning infrastructure with GitHub Actions.

## Workflow

The workflow in `/tmp/workspace/asnd/net_utils/.github/workflows/terraform.yml`:

- runs `terraform init`, `terraform fmt -check`, and `terraform plan -input=false` for pull requests
- runs `terraform apply -auto-approve -input=false` for pushes to `main`

## Required setup

Before the workflow can succeed, add:

1. a root `main.tf` that configures the Terraform Cloud remote backend and your resources
2. a repository secret named `TF_API_TOKEN` with a Terraform Cloud user API token
