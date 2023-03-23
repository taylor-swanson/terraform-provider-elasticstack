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

func TestAccResourceAgentPolicy(t *testing.T) {
	policyName := sdkacctest.RandStringFromCharSet(22, sdkacctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(t) },
		CheckDestroy:             checkResourceAgentPolicyDestroy,
		ProtoV5ProviderFactories: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAgentPolicyCreate(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_agent_policy.test_policy", "name", fmt.Sprintf("Policy %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_agent_policy.test_policy", "namespace", "default"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_agent_policy.test_policy", "description", "Test Agent Policy"),
				),
			},
			{
				Config: testAccResourceAgentPolicyUpdate(policyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_agent_policy.test_policy", "name", fmt.Sprintf("Updated Policy %s", policyName)),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_agent_policy.test_policy", "namespace", "default"),
					resource.TestCheckResourceAttr("elasticstack_kibana_fleet_agent_policy.test_policy", "description", "This policy was updated"),
				),
			},
		},
	})
}

func testAccResourceAgentPolicyCreate(id string) string {
	return fmt.Sprintf(`
provider "elasticstack" {
  elasticsearch {
    insecure = true
  }
  kibana {
    insecure = true
  }
}

resource "elasticstack_kibana_fleet_agent_policy" "test_policy" {
  name        = "%s"
  namespace   = "default"
  description = "Test Agent Policy"
  monitoring_enabled = [
    "metrics",
    "logs"
  ]
}
`, fmt.Sprintf("Policy %s", id))
}

func testAccResourceAgentPolicyUpdate(id string) string {
	return fmt.Sprintf(`
provider "elasticstack" {
  elasticsearch {
    insecure = true
  }
  kibana {
    insecure = true
  }
}

resource "elasticstack_kibana_fleet_agent_policy" "test_policy" {
  name        = "%s"
  namespace   = "default"
  description = "This policy was updated"
  monitoring_enabled = [
    "metrics",
    "logs"
  ]
}
`, fmt.Sprintf("Updated Policy %s", id))
}

func checkResourceAgentPolicyDestroy(s *terraform.State) error {
	client, err := clients.NewAcceptanceTestingClient()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "elasticstack_kibana_fleet_agent_policy" {
			continue
		}

		kibanaClient, err := client.GetKibanaClient()
		if err != nil {
			return err
		}

		res, _ := kibanaClient.Client.NewRequest().Get("/api/fleet/agent_policies/" + rs.Primary.ID)
		if res.StatusCode() != 404 {
			return fmt.Errorf("agent policy (%s) still exists", rs.Primary.ID)
		}
	}
	return nil
}
