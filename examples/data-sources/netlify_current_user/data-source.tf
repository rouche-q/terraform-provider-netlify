terraform {
  required_providers {
    netlify = {
      source = "rouche-q/netlify"
    }
  }
}

provider "netlify" {}

data "netlify_current_user" "me" {}

output "test" {
  value = data.netlify_current_user.me
}
