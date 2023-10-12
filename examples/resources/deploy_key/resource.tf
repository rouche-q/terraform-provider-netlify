terraform {
  required_providers {
    netlify = {
      source = "registry.terraform.io/rouche-q/netlify"
    }
  }
}

provider "netlify" {}

resource "netlify_deploy_key" "test" {}
