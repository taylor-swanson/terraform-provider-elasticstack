package fleet_test

import (
	"fmt"
	"testing"

	"github.com/elastic/terraform-provider-elasticstack/internal/acctest"
	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourcePackagePolicy(t *testing.T) {
	policyName := sdkacctest.RandStringFromCharSet(22, sdkacctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		CheckDestroy:             checkResourcePackagePolicyDestroy,
		ProtoV5ProviderFactories: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccResourcePackagePolicyCreate(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "name", fmt.Sprintf("Package Policy %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "namespace", "default"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "description", "Test Package Policy"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "package.0.name", "winlog"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "package.0.version", "1.12.4"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "input.0.policy_template", "winlogs"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "input.0.enabled", "true"),
				),
			},
			{
				Config: testAccResourcePackagePolicyUpdate(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "name", fmt.Sprintf("Package Policy %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "namespace", "default"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "description", "This policy was updated"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "package.0.name", "winlog"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "package.0.version", "1.12.4"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "input.0.policy_template", "winlogs"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_package_policy.test", "input.0.enabled", "true"),
				),
			},
		},
	})
}

func testAccResourcePackagePolicyCreate(id string) string {
	return fmt.Sprintf(`
provider "elasticstack" {
  elasticsearch {
    insecure = true
  }
  kibana {
    insecure = true
  }
}

resource "elasticstack_kibana_fleet_agent_policy" "test" {
  name        = "Agent Policy %s"
  namespace   = "default"
  description = "Test Agent Policy"
  monitoring_enabled = [
    "metrics",
    "logs"
  ]
}

resource "elasticstack_kibana_fleet_package_policy" "test" {
  name        = "Package Policy %s"
  namespace   = "default"
  description = "Test Package Policy"
  agent_policy_id = elasticstack_kibana_fleet_agent_policy.test.id

  package {
    name    = "winlog"
    version = "1.12.4"
  }

  input {
    policy_template = "winlogs"
    type            = "winlog"
    stream {
      data_stream = "winlog"
      vars_json = jsonencode({
        channel                 = "Security"
        "data_stream.dataset"   = "winlog.security"
        preserve_original_event = false
        ignore_older            = "72h"
        language                = 0
        event_id                = "4624,4625"
        tags                    = ["security"]
      })
    }
  }
}
`, id, id)
}

func testAccResourcePackagePolicyUpdate(id string) string {
	return fmt.Sprintf(`
provider "elasticstack" {
  elasticsearch {
    insecure = true
  }
  kibana {
    insecure = true
  }
}

resource "elasticstack_kibana_fleet_agent_policy" "test" {
  name        = "Agent Policy %s"
  namespace   = "default"
  description = "Test Agent Policy"
  monitoring_enabled = [
    "metrics",
    "logs"
  ]
}

resource "elasticstack_kibana_fleet_package_policy" "test" {
  name        = "Package Policy %s"
  namespace   = "default"
  description = "This policy was updated"
  agent_policy_id = elasticstack_kibana_fleet_agent_policy.test.id

  package {
    name    = "winlog"
    version = "1.12.4"
  }

  input {
    policy_template = "winlogs"
    type            = "winlog"
    stream {
      data_stream = "winlog"
      vars_json = jsonencode({
        channel                 = "Security"
        "data_stream.dataset"   = "winlog.security"
        preserve_original_event = false
        ignore_older            = "72h"
        language                = 0
        event_id                = "4624,4625"
        tags                    = ["security"]
      })
    }
  }
}
`, id, id)
}

func checkResourcePackagePolicyDestroy(s *terraform.State) error {
	client, err := clients.NewAcceptanceTestingClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_kibana_fleet_package_policy" {
			continue
		}

		kibanaClient, err := client.GetKibanaClient()
		if err != nil {
			return err
		}

		res, _ := kibanaClient.Client.NewRequest().Get("/api/fleet/package_policies/" + rs.Primary.ID)
		if res.StatusCode() != 404 {
			return fmt.Errorf("package policy (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}
