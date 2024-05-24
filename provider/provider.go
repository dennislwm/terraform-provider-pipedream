package provider

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "terraform-provider-pipedream/resources"
)

func Provider() *schema.Provider {
    return &schema.Provider{
        ResourcesMap: map[string]*schema.Resource{
            "pipedream_workflow": resources.ResourceWorkflow(),
        },
    }
}