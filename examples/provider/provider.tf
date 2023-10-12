terraform {
  required_providers {
    netlify = {
      source = "rouche-q/netlify"
    }
  }
}

provider "netlify" {
  personal_token = "YOUR_NETLIFY_TOKEN_HERE"
}
