# Forward Email Terraform Provider

Terraform provider for email forwarding configuration.

## Usage Example

```terraform
provider "forwardemail" {
  api_key = "your_api_key"
}

data forwardemail_account "account" {}

output "account_email" {
  value = data.forwardemail_account.account.email
}
```
