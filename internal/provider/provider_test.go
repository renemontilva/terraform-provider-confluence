package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"confluence": providerserver.NewProtocol6WithError(New()),
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CONFLUENCE_HOST"); v == "" {
		t.Fatal("CONFLUENCE_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("CONFLUENCE_USER"); v == "" {
		t.Fatal("CONFLUENCE_USER must be set for acceptance tests")
	}
	if v := os.Getenv("CONFLUENCE_TOKEN"); v == "" {
		t.Fatal("CONFLUENCE_TOKEN must be set for acceptance tests")
	}
}
