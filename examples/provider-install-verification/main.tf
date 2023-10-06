terraform {
  required_providers {
    netlify = {
      source = "registry.terraform.io/rouche-q/netlify"
    }
  }
}

provider "netlify" {
  personal_token = "YOUR_NETLIFY_TOKEN_HERE"
}

data "netlify_site" "tests" {
  id = "YOUR_NETLIFY_SITE_ID_HERE"
}

output "test" {
  value = data.netlify_site.tests
}

