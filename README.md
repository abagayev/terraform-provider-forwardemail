# Forward Email Terraform Provider

![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/abagayev/terraform-provider-forwardemail/ci.yml)
[![CodeQL](https://github.com/MichaelCurrin/badge-generator/workflows/CodeQL/badge.svg)](https://github.com/abagayev/terraform-provider-forwardemail/actions?query=workflow%3ACodeQL "Code quality workflow status")
[![codecov](https://codecov.io/gh/abagayev/terraform-provider-forwardemail/graph/badge.svg?token=R7pfHzQx3k)](https://codecov.io/gh/abagayev/terraform-provider-forwardemail)

Terraform provider for email forwarding configuration.

## Usage Example

```terraform
provider "forwardemail" {
  api_key = "your_api_key"
}

data forwardemail_account "account" {}

resource forwardemail_domain domain {
  name = "test.com"
}
```
