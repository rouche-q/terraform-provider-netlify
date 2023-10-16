terraform {
  required_providers {
    netlify = {
      source = "rouche-q/netlify"
    }
  }
}

provider "netlify" {}

data "netlify_site" "test" {
  id = "NETLIFY_SITE_ID"
}

output "test" {
  value = data.netlify_site.test
}
