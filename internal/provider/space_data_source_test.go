package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSpaceDataSourceBasic(t *testing.T) {
	resource.Test(t,
		resource.TestCase{
			PreCheck: func() {
				testAccPreCheck(t)
			},
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: testAccSpaceDataSourceConfigBasic(),
					Check: resource.ComposeAggregateTestCheckFunc(
						resource.TestCheckResourceAttr("data.confluence_space.test", "key", "DEVOPS"),
						resource.TestCheckResourceAttr("data.confluence_space.test", "type", "global"),
						resource.TestCheckResourceAttr("data.confluence_space.test", "name", "devops"),
						resource.TestCheckResourceAttr("data.confluence_space.test", "status", "current"),
					),
				},
			},
		},
	)
}

func testAccSpaceDataSourceConfigBasic() string {
	return fmt.Sprint(`
	data "confluence_space" "test" {
		key = "devops"
	}
	`)
}
