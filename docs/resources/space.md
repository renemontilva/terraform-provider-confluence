---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "confluence_space Resource - terraform-provider-confluence"
subcategory: ""
description: |-
  Creates spaces, space is a container for organizing and grouping related pages of content.
          Spaces can be used to separate content by project, team, department, or other criteria.
      A Confluence space typically includes a set of pages that can be organized into a hierarchical structure, with parent and child pages.
---

# confluence_space (Resource)

Creates spaces, space is a container for organizing and grouping related pages of content.
		Spaces can be used to separate content by project, team, department, or other criteria.
		
		A Confluence space typically includes a set of pages that can be organized into a hierarchical structure, with parent and child pages.

## Example Usage

```terraform
resource "confluence_space" "example" {
  key  = "devops"
  name = "devops"
  type = "global"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) The key of the space to be returned, which is the space name e.g: key="DEVOPS"
- `name` (String) The name of the space.

### Optional

- `description` (String) The description of the new/updated space.

### Read-Only

- `id` (Number) Space identifier number.


