terraform {
  required_providers {
    netlify = {
      source = "rouche-q/netlify"
    }
  }
}

provider "netlify" {}

resource "netlify_deploy_key" "test" {}
