# `terraform-provider-slack`

A Slack Terraform provider specializing in emoji transport.

### Where do I get a token?

Log in to Slack via a web browser. Get your token from `boot_data.api_token` in the developer tools attached to that window (right click, inspect element, console tab).

### How do I use it?

1. Build the `terraform-provider-slack` binary by running `go get github.com/christianalexander/terraform-provider-slack` or cloning and running `go build -o terraform-provider-slack .` (the executable name is important for terraform to recognize it).
1. Install the provider in the `terraform.d` directory next to your terraform repository [as documented here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).
