---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "netlify Provider"
subcategory: ""
description: |-
  
---

# netlify Provider



## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `personal_token` (String) Netlify personal token for the Netlify API. May aslo be provided via NETLIFY_PERSONAL_TOKEN env variable
