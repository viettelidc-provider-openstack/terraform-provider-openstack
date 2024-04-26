package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func TestAccNetworkingV2Trunk_nosubports(t *testing.T) {
	var port1 ports.Port
	var trunk1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccSkipReleasesBelow(t, "stable/yoga")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2TrunkNoSubports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.parent_port_1", &port1),
					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{}, &trunk1),
					resource.TestCheckResourceAttr(
						"viettelidc_networking_trunk_v2.trunk_1", "name", "trunk_1"),
					resource.TestCheckResourceAttr(
						"viettelidc_networking_trunk_v2.trunk_1", "description", "trunk_1 description"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Trunk_subports(t *testing.T) {
	var parentPort1, subport1, subport2 ports.Port
	var trunk1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccSkipReleasesBelow(t, "stable/yoga")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2TrunkSubports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.parent_port_1", &parentPort1),
					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_1", &subport1),
					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_2", &subport2),
					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{"viettelidc_networking_port_v2.subport_1", "viettelidc_networking_port_v2.subport_2"}, &trunk1, &subport1, &subport2),
				),
			},
		},
	})
}

func TestAccNetworkingV2Trunk_tags(t *testing.T) {
	var parentPort1 ports.Port
	var trunk1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccSkipReleasesBelow(t, "stable/yoga")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2TrunkTags1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.parent_port_1", &parentPort1),
					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{}, &trunk1),
					testAccCheckNetworkingV2Tags("viettelidc_networking_trunk_v2.trunk_1", []string{"a", "b", "c"}),
				),
			},
			{
				Config: testAccNetworkingV2TrunkTags2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.parent_port_1", &parentPort1),
					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{}, &trunk1),
					testAccCheckNetworkingV2Tags("viettelidc_networking_trunk_v2.trunk_1", []string{"c", "d", "e"}),
				),
			},
		},
	})
}

// NOTE: this test is flacky and can fail with the following error:
// X-Openstack-Request-Id: req-1f854d77-414b-4826-9f24-95cac6cda10c
// 2021/09/18 11:24:00 [DEBUG] OpenStack Response Body: {
//   "NeutronError": {
//     "detail": "",
//     "message": "Unable to complete operation on port f81d3dcd-6069-4c60-8837-19d6f4abf52e for network 30dddb37-e5dd-4b71-91c3-bc1e9c3066cb. Port already has an attached device 4e227e3c-3231-4d63
// -b58f-eb7731b4480a.",
//     "type": "PortInUse"
//   }
// }
// 2021/09/18 11:24:01 [WARN] Got error running Terraform: exit status 1
//
// Error: Error updating viettelidc_networking_trunk_v2 4e227e3c-3231-4d63-b58f-eb7731b4480a subports: Expected HTTP response code [200] when accessing [PUT http://192.168.0.118:9696/v2.0/trunks/4e2
// 27e3c-3231-4d63-b58f-eb7731b4480a/add_subports], but got 409 instead
// {"NeutronError": {"type": "PortInUse", "message": "Unable to complete operation on port f81d3dcd-6069-4c60-8837-19d6f4abf52e for network 30dddb37-e5dd-4b71-91c3-bc1e9c3066cb. Port already has an
// attached device 4e227e3c-3231-4d63-b58f-eb7731b4480a.", "detail": ""}}
//
//   with viettelidc_networking_trunk_v2.trunk_1,
//   on terraform_plugin_test.tf line 44, in resource "viettelidc_networking_trunk_v2" "trunk_1":
//   44: resource "viettelidc_networking_trunk_v2" "trunk_1" {
//
//     TestAccNetworkingV2Trunk_trunkUpdateSubports: resource_viettelidc_networking_trunk_v2_test.go:103: Step 2/4 error: Error running apply: exit status 1
//
//         Error: Error updating viettelidc_networking_trunk_v2 4e227e3c-3231-4d63-b58f-eb7731b4480a subports: Expected HTTP response code [200] when accessing [PUT http://192.168.0.118:9696/v2.0/tr
// unks/4e227e3c-3231-4d63-b58f-eb7731b4480a/add_subports], but got 409 instead
//         {"NeutronError": {"type": "PortInUse", "message": "Unable to complete operation on port f81d3dcd-6069-4c60-8837-19d6f4abf52e for network 30dddb37-e5dd-4b71-91c3-bc1e9c3066cb. Port alread
// y has an attached device 4e227e3c-3231-4d63-b58f-eb7731b4480a.", "detail": ""}}
//
//           with viettelidc_networking_trunk_v2.trunk_1,
//           on terraform_plugin_test.tf line 44, in resource "viettelidc_networking_trunk_v2" "trunk_1":
//           44: resource "viettelidc_networking_trunk_v2" "trunk_1" {
//func TestAccNetworkingV2Trunk_trunkUpdateSubports(t *testing.T) {
//	var parentPort1, subport1, subport2, subport3, subport4 ports.Port
//	var trunk1 trunks.Trunk
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//			testAccPreCheckNonAdminOnly(t)
//		},
//		ProviderFactories: testAccProviders,
//		CheckDestroy:      testAccCheckNetworkingV2TrunkDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccNetworkingV2TrunkUpdateSubports1,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.parent_port_1", &parentPort1),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_1", &subport1),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_2", &subport2),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_3", &subport3),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_4", &subport4),
//					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{"viettelidc_networking_port_v2.subport_1", "viettelidc_networking_port_v2.subport_2"}, &trunk1, &subport1, &subport2),
//					resource.TestCheckResourceAttr(
//						"viettelidc_networking_trunk_v2.trunk_1", "description", "trunk_1 description"),
//				),
//			},
//			{
//				Config: testAccNetworkingV2TrunkUpdateSubports2,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.parent_port_1", &parentPort1),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_1", &subport1),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_2", &subport2),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_3", &subport3),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_4", &subport4),
//					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{"viettelidc_networking_port_v2.subport_1", "viettelidc_networking_port_v2.subport_3", "viettelidc_networking_port_v2.subport_4"}, &trunk1, &subport1, &subport3, &subport4),
//					resource.TestCheckResourceAttr(
//						"viettelidc_networking_trunk_v2.trunk_1", "description", ""),
//				),
//			},
//			{
//				Config: testAccNetworkingV2TrunkUpdateSubports3,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.parent_port_1", &parentPort1),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_1", &subport1),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_2", &subport2),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_3", &subport3),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_4", &subport4),
//					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{"viettelidc_networking_port_v2.subport_1", "viettelidc_networking_port_v2.subport_3", "viettelidc_networking_port_v2.subport_4"}, &trunk1, &subport1, &subport3, &subport4),
//					resource.TestCheckResourceAttr(
//						"viettelidc_networking_trunk_v2.trunk_1", "description", ""),
//				),
//			},
//			{
//				Config: testAccNetworkingV2TrunkUpdateSubports4,
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.parent_port_1", &parentPort1),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_1", &subport1),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_2", &subport2),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_3", &subport3),
//					testAccCheckNetworkingV2PortExists("viettelidc_networking_port_v2.subport_4", &subport4),
//					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{}, &trunk1),
//					resource.TestCheckResourceAttr(
//						"viettelidc_networking_trunk_v2.trunk_1", "description", "trunk_1 updated description"),
//				),
//			},
//		},
//	})
//}

func TestAccNetworkingV2Trunk_Instance(t *testing.T) {
	var instance1 servers.Server
	var parentPort1, subport1 ports.Port
	var trunk1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccSkipReleasesBelow(t, "stable/yoga")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2TrunkComputeInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("viettelidc_compute_instance_v2.instance_1", &instance1),
					testAccCheckNetworkingV2PortExists(
						"viettelidc_networking_port_v2.parent_port_1", &parentPort1),
					testAccCheckNetworkingV2PortExists(
						"viettelidc_networking_port_v2.subport_1", &subport1),
					testAccCheckNetworkingV2TrunkExists("viettelidc_networking_trunk_v2.trunk_1", []string{"viettelidc_networking_port_v2.subport_1"}, &trunk1, &subport1),
					resource.TestCheckResourceAttrPtr(
						"viettelidc_compute_instance_v2.instance_1", "network.0.port", &trunk1.PortID),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2TrunkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "viettelidc_networking_trunk_v2" {
			continue
		}

		_, err := trunks.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Trunk still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2TrunkExists(n string, subportResourceNames []string, trunk *trunks.Trunk, subports ...*ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Trunk not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Trunk ID is not set")
		}

		var subportResources map[string]bool
		if len(subports) > 0 {
			if len(subportResourceNames) != len(subports) {
				return fmt.Errorf("Amount of subport resource names and subports do not match")
			}

			subportResources = make(map[string]bool)
			for i, subport := range subports {
				if subportResource, ok := s.RootModule().Resources[subportResourceNames[i]]; ok {
					subportResources[subportResource.Primary.ID] = true
				} else {
					return fmt.Errorf("Subport not found: %s", subport.ID)
				}
			}
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := trunks.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if len(found.Subports) != len(subports) {
			return fmt.Errorf("The amount of retrieved trunk subports and trunk subports to check does not match")
		}

		if len(subports) > 0 {
			for _, subport := range found.Subports {
				if _, ok := subportResources[subport.PortID]; !ok {
					return fmt.Errorf("Trunk Subport not found: %s", subport.PortID)
				}
			}
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Trunk not found")
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Trunk name does not match")
		}

		*trunk = *found

		return nil
	}
}

const testAccNetworkingV2TrunkNoSubports = `
resource "viettelidc_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "viettelidc_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  description = "trunk_1 description"
  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"
}
`

const testAccNetworkingV2TrunkSubports = `
resource "viettelidc_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "viettelidc_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_port_v2" "subport_1" {
  name = "subport_1"
  admin_state_up = "true"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_port_v2" "subport_2" {
  name = "subport_2"
  admin_state_up = "true"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  description = "trunk_1 description"
  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"

  sub_port {
	  port_id = "${viettelidc_networking_port_v2.subport_1.id}"
	  segmentation_id = 1
	  segmentation_type = "vlan"
  }

  sub_port {
	  port_id = "${viettelidc_networking_port_v2.subport_2.id}"
	  segmentation_id = 2
	  segmentation_type = "vlan"
  }
}
`

//const testAccNetworkingV2TrunkUpdateSubports1 = `
//resource "viettelidc_networking_network_v2" "network_1" {
//  name = "network_1"
//  admin_state_up = "true"
//}
//
//resource "viettelidc_networking_subnet_v2" "subnet_1" {
//  name = "subnet_1"
//  cidr = "192.168.199.0/24"
//  ip_version = 4
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "parent_port_1" {
//  name = "port_1"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_1" {
//  name = "subport_1"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_2" {
//  name = "subport_2"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_3" {
//  name = "subport_3"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_4" {
//  name = "subport_4"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_trunk_v2" "trunk_1" {
//  name = "trunk_1"
//  description = "trunk_1 description"
//  admin_state_up = "true"
//  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"
//
//  sub_port {
//	  port_id = "${viettelidc_networking_port_v2.subport_1.id}"
//	  segmentation_id = 1
//	  segmentation_type = "vlan"
//  }
//
//  sub_port {
//	  port_id = "${viettelidc_networking_port_v2.subport_2.id}"
//	  segmentation_id = 2
//	  segmentation_type = "vlan"
//  }
//}
//`
//
//const testAccNetworkingV2TrunkUpdateSubports2 = `
//resource "viettelidc_networking_network_v2" "network_1" {
//  name = "network_1"
//  admin_state_up = "true"
//}
//
//resource "viettelidc_networking_subnet_v2" "subnet_1" {
//  name = "subnet_1"
//  cidr = "192.168.199.0/24"
//  ip_version = 4
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "parent_port_1" {
//  name = "port_1"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_1" {
//  name = "subport_1"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_2" {
//  name = "subport_2"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_3" {
//  name = "subport_3"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_4" {
//  name = "subport_4"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_trunk_v2" "trunk_1" {
//  name = "update_trunk_1"
//  admin_state_up = "true"
//  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"
//
//  sub_port {
//	  port_id = "${viettelidc_networking_port_v2.subport_1.id}"
//	  segmentation_id = 1
//	  segmentation_type = "vlan"
//  }
//
//  sub_port {
//	  port_id = "${viettelidc_networking_port_v2.subport_3.id}"
//	  segmentation_id = 3
//	  segmentation_type = "vlan"
//  }
//
//  sub_port {
//	  port_id = "${viettelidc_networking_port_v2.subport_4.id}"
//	  segmentation_id = 4
//	  segmentation_type = "vlan"
//  }
//}
//`
//
//const testAccNetworkingV2TrunkUpdateSubports3 = `
//resource "viettelidc_networking_network_v2" "network_1" {
//  name = "network_1"
//  admin_state_up = "true"
//}
//
//resource "viettelidc_networking_subnet_v2" "subnet_1" {
//  name = "subnet_1"
//  cidr = "192.168.199.0/24"
//  ip_version = 4
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "parent_port_1" {
//  name = "port_1"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_1" {
//  name = "subport_1"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_2" {
//  name = "subport_2"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_3" {
//  name = "subport_3"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_4" {
//  name = "subport_4"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_trunk_v2" "trunk_1" {
//  name = "trunk_1"
//  description = ""
//  admin_state_up = "true"
//  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"
//
//  sub_port {
//	  port_id = "${viettelidc_networking_port_v2.subport_1.id}"
//	  segmentation_id = 1
//	  segmentation_type = "vlan"
//  }
//
//  sub_port {
//	  port_id = "${viettelidc_networking_port_v2.subport_3.id}"
//	  segmentation_id = 3
//	  segmentation_type = "vlan"
//  }
//
//  sub_port {
//	  port_id = "${viettelidc_networking_port_v2.subport_4.id}"
//	  segmentation_id = 4
//	  segmentation_type = "vlan"
//  }
//}
//`
//
//const testAccNetworkingV2TrunkUpdateSubports4 = `
//resource "viettelidc_networking_network_v2" "network_1" {
//  name = "network_1"
//  admin_state_up = "true"
//}
//
//resource "viettelidc_networking_subnet_v2" "subnet_1" {
//  name = "subnet_1"
//  cidr = "192.168.199.0/24"
//  ip_version = 4
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "parent_port_1" {
//  name = "port_1"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_1" {
//  name = "subport_1"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_2" {
//  name = "subport_2"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_3" {
//  name = "subport_3"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_port_v2" "subport_4" {
//  name = "subport_4"
//  admin_state_up = "true"
//  network_id = "${viettelidc_networking_network_v2.network_1.id}"
//}
//
//resource "viettelidc_networking_trunk_v2" "trunk_1" {
//  name = "trunk_1"
//  description = "trunk_1 updated description"
//  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"
//  admin_state_up = "true"
//}
//`

const testAccNetworkingV2TrunkComputeInstance = `
resource "viettelidc_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "viettelidc_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "viettelidc_networking_port_v2" "parent_port_1" {
  depends_on = [
    "viettelidc_networking_subnet_v2.subnet_1",
  ]

  name = "parent_port_1"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

resource "viettelidc_networking_port_v2" "subport_1" {
  depends_on = [
    "viettelidc_networking_subnet_v2.subnet_1",
  ]

  name = "subport_1"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

resource "viettelidc_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  admin_state_up = "true"
  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"

  sub_port {
	  port_id = "${viettelidc_networking_port_v2.subport_1.id}"
	  segmentation_id = 1
	  segmentation_type = "vlan"
  }
}

resource "viettelidc_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]

  network {
    port = "${viettelidc_networking_trunk_v2.trunk_1.port_id}"
  }
}
`

const testAccNetworkingV2TrunkTags1 = `
resource "viettelidc_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "viettelidc_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"

  tags = ["a", "b", "c"]
}
`

const testAccNetworkingV2TrunkTags2 = `
resource "viettelidc_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "viettelidc_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${viettelidc_networking_network_v2.network_1.id}"
}

resource "viettelidc_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  port_id = "${viettelidc_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"

  tags = ["c", "d", "e"]
}
`
