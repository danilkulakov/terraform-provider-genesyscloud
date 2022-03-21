---
page_title: "genesyscloud_webdeployments_deployment Resource - terraform-provider-genesyscloud"
subcategory: ""
description: |-
  Genesys Cloud Web Deployment
---
# genesyscloud_webdeployments_deployment (Resource)

Genesys Cloud Web Deployment

## API Usage
The following Genesys Cloud APIs are used by this resource. Ensure your OAuth Client has been granted the necessary scopes and permissions to perform these operations:

* [GET /api/v2/webdeployments/deployments](https://developer.dev-genesys.cloud/api/rest/v2/webdeployments/#get-api-v2-webdeployments-deployments)
* [POST /api/v2/webdeployments/deployments](https://developer.dev-genesys.cloud/api/rest/v2/webdeployments/#post-api-v2-webdeployments-deployments)
* [DELETE /api/v2/webdeployments/deployments/{deploymentId}](https://developer.dev-genesys.cloud/api/rest/v2/webdeployments/#delete-api-v2-webdeployments-deployments--deploymentId-)
* [GET /api/v2/webdeployments/deployments/{deploymentId}](https://developer.dev-genesys.cloud/api/rest/v2/webdeployments/#get-api-v2-webdeployments-deployments--deploymentId-)
* [PUT /api/v2/webdeployments/deployments/{deploymentId}](https://developer.dev-genesys.cloud/api/rest/v2/webdeployments/#put-api-v2-webdeployments-deployments--deploymentId-)

## Example Usage

```terraform
data "genesyscloud_flow" "incomingMessageFlow" {
  name = "Incoming Message Flow"
}

data "genesyscloud_webdeployments_configuration" "exampleConfiguration" {
  name = "Example Web Deployment Configuration"
}

resource "genesyscloud_webdeployments_deployment" "exampleDeployment" {
  name            = "Example Web Deployment"
  description     = "This is an example of a web deployment"
  allowed_domains = ["genesys.com"]
  flow_id         = data.genesyscloud_flow.incomingMessageFlow.id
  configuration {
    id      = data.genesyscloud_webdeployments_configuration.exampleConfiguration.id
    version = data.genesyscloud_webdeployments_configuration.exampleConfiguration.version
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **configuration** (Block List, Min: 1, Max: 1) The published configuration version used by this deployment (see [below for nested schema](#nestedblock--configuration))
- **name** (String) Deployment name

### Optional

- **allow_all_domains** (Boolean) Whether all domains are allowed or not. allowedDomains must be empty when this is true. Defaults to `false`.
- **allowed_domains** (List of String) The list of domains that are approved to use this deployment; the list will be added to CORS headers for ease of web use.
- **description** (String) Deployment description
- **flow_id** (String) A reference to the inboundshortmessage flow used by this deployment.
- **id** (String) The ID of this resource.
- **status** (String) The current status of the deployment. Valid values: Pending, Active, Inactive, Error, Deleting.

<a id="nestedblock--configuration"></a>
### Nested Schema for `configuration`

Required:

- **id** (String) The ID of this resource.

Optional:

- **version** (String)
