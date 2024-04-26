package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSiteConnectionV2_importBasic(t *testing.T) {
	resourceName := "viettelidc_vpnaas_site_connection_v2.conn_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSiteConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConnectionV2Basic(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
