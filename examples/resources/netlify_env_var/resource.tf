terraform {
  required_providers {
    netlify = {
      source = "rouche-q/netlify"
    }
  }
}

provider "netlify" {}

resource "netlify_deploy_key" "test" {}

data "netlify_current_user" "me" {}

output "key" {
  value = resource.netlify_deploy_key.test
}

resource "netlify_site" "test" {
  repository = {
    provider      = "github"
    repo_path     = "USER/repo"
    repo_branch   = "main"
    deploy_key_id = resource.netlify_deploy_key.test.id
    cmd           = "npm run build"
    dir           = "build"
  }
}

resource "netlify_env_var" "test_env" {
  account_slug = data.netlify_current_user.me.slug
  site_id      = resource.netlify_site.test.id
  key          = "test"
}
