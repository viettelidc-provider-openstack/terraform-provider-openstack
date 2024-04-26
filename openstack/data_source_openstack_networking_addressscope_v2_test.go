package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOpenStackNetworkingAddressScopeV2DataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSourceAddressscope,
			},
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.viettelidc_networking_addressscope_v2.addressscope_1"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "name", "addressscope_1"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "ip_version", "4"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "shared", "false"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingAddressScopeV2DataSource_ipversion(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSourceAddressscope,
			},
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSourceIPVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.viettelidc_networking_addressscope_v2.addressscope_1"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "name", "addressscope_1"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "ip_version", "4"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "shared", "false"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingAddressScopeV2DataSource_shared(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSourceAddressscope,
			},
			{
				Config: testAccOpenStackNetworkingAddressScopeV2DataSourceShared(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.viettelidc_networking_addressscope_v2.addressscope_1"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "name", "addressscope_1"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "ip_version", "4"),
					resource.TestCheckResourceAttr("viettelidc_networking_addressscope_v2.addressscope_1", "shared", "false"),
				),
			},
		},
	})
}

const testAccOpenStackNetworkingAddressScopeV2DataSourceAddressscope = `
resource "viettelidc_networking_addressscope_v2" "addressscope_1" {
  name       = "addressscope_1"
  ip_version = 4
  shared     = false
}`

func testAccOpenStackNetworkingAddressScopeV2DataSourceName() string {
	return fmt.Sprintf(`
%s

data "viettelidc_networking_addressscope_v2" "addressscope_1" {
  name = "${viettelidc_networking_addressscope_v2.addressscope_1.name}"
}
`, testAccOpenStackNetworkingAddressScopeV2DataSourceAddressscope)
}

func testAccOpenStackNetworkingAddressScopeV2DataSourceIPVersion() string {
	return fmt.Sprintf(`
%s

data "viettelidc_networking_addressscope_v2" "addressscope_1" {
  ip_version = "${viettelidc_networking_addressscope_v2.addressscope_1.ip_version}"
}
`, testAccOpenStackNetworkingAddressScopeV2DataSourceAddressscope)
}

func testAccOpenStackNetworkingAddressScopeV2DataSourceShared() string {
	return fmt.Sprintf(`
%s

data "viettelidc_networking_addressscope_v2" "addressscope_1" {
  shared = "${viettelidc_networking_addressscope_v2.addressscope_1.shared}"
}
`, testAccOpenStackNetworkingAddressScopeV2DataSourceAddressscope)
}
