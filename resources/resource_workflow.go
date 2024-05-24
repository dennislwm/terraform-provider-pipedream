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