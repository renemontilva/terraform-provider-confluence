package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSpaceResourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Create and Read Space
				Config: testAccSpaceResourceConfigBasic("TestSpace", "test description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("confluence_space.test", "name", "TestSpace"),
					resource.TestCheckResourceAttr("confluence_space.test", "description", "test description"),
				),
			},
			{
				ResourceName:      "confluence_space.test",
				ImportState:       true,
				ImportStateId:     "TestSpace",
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"description",
				},
			},
			{
				// Update and Read Space
				Config: testAccSpaceResourceConfigBasic("TestSpaceUpdate", "test description update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("confluence_space.test", "name", "TestSpaceUpdate"),
					resource.TestCheckResourceAttr("confluence_space.test", "description", "test description update"),
				),
			},
		},
	})
}

func testAccSpaceResourceConfigBasic(name, description string) string {
	return fmt.Sprintf(`
	resource "confluence_space" "test" {
		key = "TestSpace"
		name = "%s"
		description = "%s"
	}
	`, name, description)
}
