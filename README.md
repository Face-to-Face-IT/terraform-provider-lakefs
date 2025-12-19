# Terraform Provider for LakeFS

[![Registry](https://img.shields.io/badge/registry-Face-to-Face-IT%2Flakefs-blue)](https://registry.terraform.io/providers/Face-to-Face-IT/lakefs/latest)

This Terraform provider allows you to manage [LakeFS](https://lakefs.io/) resources using Infrastructure as Code.

LakeFS is an open-source data version control system for data lakes. It provides Git-like operations such as branching, committing, and merging for your data.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for building from source)

## Using the Provider

```hcl
terraform {
  required_providers {
    lakefs = {
      source  = "Face-to-Face-IT/lakefs"
      version = "~> 0.1"
    }
  }
}

provider "lakefs" {
  endpoint          = "http://localhost:8000/api/v1"
  access_key_id     = "AKIAIOSFODNN7EXAMPLE"
  secret_access_key = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
}
```

See the [provider documentation](https://registry.terraform.io/providers/Face-to-Face-IT/lakefs/latest/docs) for full configuration details.

## Features

### Resources
- `lakefs_repository` - Manage repositories
- `lakefs_branch` - Manage branches
- `lakefs_tag` - Manage tags
- `lakefs_branch_protection` - Manage branch protection rules
- `lakefs_user` - Manage users
- `lakefs_group` - Manage groups
- `lakefs_policy` - Manage policies
- `lakefs_group_membership` - Manage group memberships
- `lakefs_user_policy_attachment` - Attach policies to users
- `lakefs_group_policy_attachment` - Attach policies to groups
- `lakefs_user_credentials` - Manage user credentials

### Data Sources
- `lakefs_repository` - Query repository info
- `lakefs_branch` - Query branch info
- `lakefs_commit` - Query commit info
- `lakefs_current_user` - Query authenticated user
- `lakefs_user` - Query user info
- `lakefs_group` - Query group info
- `lakefs_policy` - Query policy info

## Developing the Provider

### Building

```shell
go build -o terraform-provider-lakefs
```

### Testing

```shell
# Start test infrastructure
make testacc-up

# Run acceptance tests
make testacc-local

# Stop test infrastructure
make testacc-down
```

### Generating Documentation

```shell
go generate ./...
```

## License

MPL-2.0 - see [LICENSE](LICENSE)
