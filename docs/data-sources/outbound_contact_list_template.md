---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "genesyscloud_outbound_contact_list_template Data Source - terraform-provider-genesyscloud"
subcategory: ""
description: |-
  Data source for Genesys Cloud Outbound Contact Lists Templates. Select a contact list template by name.
---

# genesyscloud_outbound_contact_list_template (Data Source)

Data source for Genesys Cloud Outbound Contact Lists Templates. Select a contact list template by name.

## Example Usage

```terraform
data "genesyscloud_outbound_contact_list_template" "contact_list_template" {
  name = "Example Contact List Template"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Contact List Template name.

### Read-Only

- `id` (String) The ID of this resource.
