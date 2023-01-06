package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccContentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccContentResourceConfig("test create"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("confluence_content.test", "body", "<h1>Terraform Acc test create</h1><p>paragraph</p>"),
					resource.TestCheckResourceAttr("confluence_content.test", "space", "DEVOPS"),
					resource.TestCheckResourceAttr("confluence_content.test", "title", "test create"),
				),
			},
			// ImportState testing
			{
				ResourceName: "confluence_content.test",
				ImportState:  true,
				// FIXME: Because the GetContentById does not retrieve body field, rease why is set to false
				ImportStateVerify: false,
			},
			// Update and Read testing
			{
				Config: testAccContentResourceConfig("test update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("confluence_content.test", "body", "<h1>Terraform Acc test update</h1><p>paragraph</p>"),
					resource.TestCheckResourceAttr("confluence_content.test", "space", "DEVOPS"),
					resource.TestCheckResourceAttr("confluence_content.test", "title", "test update"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccContentResourceConfig(title string) string {
	return fmt.Sprintf(`
resource "confluence_content" "test" {
  type = "page" 
  title = "%s"
  space = "DEVOPS"
  body = "<h1>Terraform Acc %s</h1><p>paragraph</p>"
}`, title, title)
}
