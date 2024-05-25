# terraform-provider-pipedream

---
# 1. Introduction
## 1.1. Purpose

This document describes the `terrafom-provider-pipedream` custom Terraform provider for Pipedream, an SaaS service that provides workflow automation similar to Zapier, and n8n.

## 1.2. Audience

The audience for this document includes:

* Developer who will develop the provider, run unit tests, configure build tools and write user documentation.

* DevSecOps Engineer who will shape the workflow, create and maintain continuous integration and deployment (CI/CD) pipelines, and write playbooks and runbooks.

---
# 2. System Overview
## 2.1. Benefits and Values

1. Currently, creating a Pipedream workflow requires the User to navigate the Pipedream user interface (UI) and perform click operations, which may be inefficient and error prone.

2. The `terrafom-provider-pipedream` allows the User to create and maintain Terraform configuration files that define a Pipedream workflow as code, hence reducing the error rate, while increasing the reusability and adding version control for a configuration as code.

## 2.2. Workflow

This project uses several methods and products to optimize your workflow.
- Use a version control system (**GitHub**) to track your changes and collaborate with others.
- Use a cloud LLM (**ChatGPT**) to facilitate shaping and writing playbook and runbooks.
- Use a Python LLM-enabled CLI (**Aider.chat**) to facilitate coding.
- Use a CLI (**restish**) to interact with an API endpoint.
- Use a compiler (**Go**) to compile the source code into a distributable file.
- Use a build tool (**Makefile**) to automate your build tasks.
- Use a containerization platform (**Docker**) to run your application in any environment.

---
# 3. User Personas
## 3.1 RACI Matrix

|            Category            |                         Activity                         | Developer | DevSecOps |
|:------------------------------:|:--------------------------------------------------------:|:---------:|:---------:|
| Installation and Configuration |    Configuring `aider.chat` in your local repository     |           |    R,A    |
|       Shaping by GPT-4o        |              Creating the project structure              |    R,A    |           |
|       Shaping by GPT-4o        |                Writing your provider code                |    R,A    |           |
|       Shaping by GPT-4o        |                  Defining your resource                  |    R,A    |           |
|       Shaping by GPT-4o        | Initializing the Go module and adding the helper modules |    R,A    |           |
|       Shaping by GPT-4o        |                  Building your provider                  |    R,A    |           |
|       Shaping by GPT-4o        |             Adding a local plugin directory              |    R,A    |           |
|       Shaping by GPT-4o        |                  Testing your provider                   |    R,A    |           |

---
# 4. Requirements
## 4.1. Local workstation

- ChatGPT Desktop for macOS
- Go 1.22.0 (`/opt/homebrew/bin/go`)
- Python 3.11.7 (`/opt/homebrew/bin/python3`)
  - `aider-chat` 0.36.0 (`python3 -m pip install`)
- Terraform 1.7.5 (`/opt/homebrew/bin/terraform`)

## 4.2. SaaS accounts

- GitHub account
- OpenAI ChatGPT Plus account

---
# 5. Installation and Configuration
## 5.1. Configuring `aider.chat` in your local repository

This runbook should be performed by the DevSecOps Engineer.

1. Open a bash terminal and navigate to your workspace > type the following command.

```sh
git clone https://github.com/dennislwm/terraform-provider-pipedream
```

2. Create an file `.env` in your editor and copy and paste the content below, replacing the `<TOKEN>` with your API token.

```txt
export OPENAI_API_KEY=<TOKEN>
```

3. Navigate to your local repository `terraform-provider-pipedream`, and type the command:

```sh
source .env
aider --version
```

> Note: You do not need to launch a virtual environment, if the package `aider-chat` was `pip` installed globally.

---
# 6. Shaping by GPT-4o
## 6.1. Creating the project structure

This runbook should be performed by the Developer.

1. Create a new directory structure for your project with the following subdirectories and files.

```sh
terraform-provider-pipedream/
|- main.go
+- provider/
   |- provider.go
+- resources/
   |- resource_workflow.go
```

## 6.2. Writing your provider code

This runbook should be performed by the Developer.

1. Create and edit the `main.go` file.

```go
package main

import (
    "github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
    "terraform-provider-pipedream/provider"
)

func main() {
    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: provider.Provider,
    })
}
```

2. Create and edit the `provider/provider.go` file.

```go
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
```

## 6.3. Defining your resource

This runbook should be performed by the Developer.

Assuming that your first endpoint is `https://api.pipedream.com/v1/workflows`.

1. Create and edit the `resources/resource_workflow.go` file.

```go
package resources

import (
    "net/http"
    "io/ioutil"
    "bytes"
    "encoding/json"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const baseURL = "https://api.pipedream.com/v1/workflows"

func ResourceWorkflow() *schema.Resource {
    return &schema.Resource{
        Create: resourceWorkflowCreate,
        Read:   resourceWorkflowRead,
        Update: resourceWorkflowUpdate,
        Delete: resourceWorkflowDelete,

        Schema: map[string]*schema.Schema{
            "name": {
                Type:     schema.TypeString,
                Required: true,
            },
            "description": {
                Type:     schema.TypeString,
                Optional: true,
            },
        },
    }
}

func resourceWorkflowCreate(d *schema.ResourceData, m interface{}) error {
    name := d.Get("name").(string)
    description := d.Get("description").(string)

    requestBody, err := json.Marshal(map[string]string{
        "name":        name,
        "description": description,
    })
    if err != nil {
        return err
    }

    resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(requestBody))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    var result map[string]interface{}
    json.Unmarshal(body, &result)
    d.SetId(result["id"].(string))

    return resourceWorkflowRead(d, m)
}

func resourceWorkflowRead(d *schema.ResourceData, m interface{}) error {
    id := d.Id()

    resp, err := http.Get(baseURL + "/" + id)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        d.SetId("")
        return nil
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    var result map[string]interface{}
    json.Unmarshal(body, &result)

    d.Set("name", result["name"])
    d.Set("description", result["description"])

    return nil
}

func resourceWorkflowUpdate(d *schema.ResourceData, m interface{}) error {
    name := d.Get("name").(string)
    description := d.Get("description").(string)
    id := d.Id()

    requestBody, err := json.Marshal(map[string]string{
        "name":        name,
        "description": description,
    })
    if err != nil {
        return err
    }

    req, err := http.NewRequest(http.MethodPut, baseURL+"/"+id, bytes.NewBuffer(requestBody))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return resourceWorkflowRead(d, m)
}

func resourceWorkflowDelete(d *schema.ResourceData, m interface{}) error {
    id := d.Id()

    req, err := http.NewRequest(http.MethodDelete, baseURL+"/"+id, nil)
    if err != nil {
        return err
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    d.SetId("")

    return nil
}
```

## 6.4. Initializing the Go module and adding the helper modules

This runbook should be performed by the Developer.

1. Open a bash terminal and run the command to initialize `terraform-provider-pipedream` module:

```sh
go mod init terraform-provider-pipedream
```

2. Run the command to add the required helper modules.

```sh
go get github.com/hashicorp/terraform-plugin-sdk/v2/plugin
go get github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema
```

## 6.5. Building your provider

This runbook should be performed by the Developer.

1. Open a bash terminal and run the command:

```sh
GOOS=darwin GOARCH=arm64 go build -o terraform-provider-pipedream
```

2. This should output a Go file `terraform-provider-pipedream` in the current folder.

## 6.6. Adding a local plugin directory

This runbook should be performed by the Developer.

1. Open a bash terminal and run the command:

```sh
mkdir -p ~/.terraform.d/plugins/localhost/local/pipedream/1.0.0/darwin_amd64
mv terraform-provider-pipedream ~/.terraform.d/plugins/localhost/local/pipedream/1.0.0/darwin_amd64/
```

## 6.7. Testing your provider

This runbook should be performed by the Developer.

1. Create a new directory `test_terraform` for your Terraform configuration.

2. Change to the directory, and create a `main.tf` file:

```hcl
terraform {
  required_providers {
    pipedream = {
      source  = "localhost/local/pipedream"
      version = "1.0.0"
    }
  }
}

provider "pipedream" {}

resource "pipedream_workflow" "example" {
  name        = "Example Workflow"
  description = "This is an example workflow."
}
```

3. Initialize Terraform and apply your configuration.

```sh
terraform init
terraform apply
```
