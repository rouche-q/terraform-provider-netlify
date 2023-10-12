terraform {
  required_providers {
    netlify = {
      source = "registry.terraform.io/rouche-q/netlify"
    }
  }
}

provider "netlify" {
}

data "netlify_site" "test" {
  id = "NETLIFY_SITE_ID"
}

output "test" {
  value = data.netlify_site.test
}
