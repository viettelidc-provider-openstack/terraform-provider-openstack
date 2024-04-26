package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2QoSMinimumBandwidthRuleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccSkipReleasesBelow(t, "stable/yoga")
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSMinimumBandwidthRuleDataSource,
			},
			{
				Config: testAccOpenStackNetworkingQoSMinimumBandwidthRuleV2DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingQoSMinimumBandwidthRuleV2DataSourceID("data.viettelidc_networking_qos_minimum_bandwidth_rule_v2.min_bw_rule_1"),
					resource.TestCheckResourceAttr(
						"data.viettelidc_networking_qos_minimum_bandwidth_rule_v2.min_bw_rule_1", "min_kbps", "3000"),
				),
			},
		},
	})
}

func testAccCheckNetworkingQoSMinimumBandwidthRuleV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find QoS minimum bw data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("QoS minimum bw data source ID not set")
		}

		return nil
	}
}

const testAccNetworkingV2QoSMinimumBandwidthRuleDataSource = `
resource "viettelidc_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "viettelidc_networking_qos_minimum_bandwidth_rule_v2" "min_bw_rule_1" {
  qos_policy_id  = "${viettelidc_networking_qos_policy_v2.qos_policy_1.id}"
  min_kbps       = 3000
}
`

func testAccOpenStackNetworkingQoSMinimumBandwidthRuleV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s
data "viettelidc_networking_qos_minimum_bandwidth_rule_v2" "min_bw_rule_1" {
  qos_policy_id = "${viettelidc_networking_qos_policy_v2.qos_policy_1.id}"
  min_kbps      = "${viettelidc_networking_qos_minimum_bandwidth_rule_v2.min_bw_rule_1.min_kbps}"
}
`, testAccNetworkingV2QoSMinimumBandwidthRuleDataSource)
}
