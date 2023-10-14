terraform {
  required_providers {
    netlify = {
      source = "rouche-q/netlify"
    }
  }
}

provider "netlify" {}

resource "netlify_deploy_key" "test" {}

output "key" {
  value = netlify_deploy_key.test
}

resource "netlify_site" "test" {
  repository = {
    provider      = "github"
    repo_path     = "USER/REPO_NAME"
    repo_branch   = "main"
    deploy_key_id = netlify_deploy_key.test.id
    cmd           = "COMMAND"
    dir           = "DIR_TO_DEPLOY"
  }
}
