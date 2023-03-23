package fleet_test

import (
	"testing"

	"github.com/elastic/terraform-provider-elasticstack/internal/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceEnrollmentToken(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		ProtoV5ProviderFactories: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceEnrollmentToken,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.elasticstack_kibana_fleet_enrollment_token.test", "id", "223b1bf8-240f-463f-8466-5062670d0754"),
				),
			},
		},
	})
}

const testAccDataSourceEnrollmentToken = `
provider "elasticstack" {
  kibana {
    insecure = true
  }
}

data "elasticstack_kibana_fleet_enrollment_token" "test" {
	key_id = "223b1bf8-240f-463f-8466-5062670d0754"
}
`
