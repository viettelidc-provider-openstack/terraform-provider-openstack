package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
)

func TestAccSiteConnectionVPNaaSV2_basic(t *testing.T) {
	var conn siteconnections.Connection
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
			t.Skip("Currently failing in GH-A")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSiteConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConnectionV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteConnectionV2Exists(
						"viettelidc_vpnaas_site_connection_v2.conn_1", &conn),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "ikepolicy_id", &conn.IKEPolicyID),
					resource.TestCheckResourceAttr("viettelidc_vpnaas_site_connection_v2.conn_1", "admin_state_up", strconv.FormatBool(conn.AdminStateUp)),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "psk", &conn.PSK),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "ipsecpolicy_id", &conn.IPSecPolicyID),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "vpnservice_id", &conn.VPNServiceID),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "local_ep_group_id", &conn.LocalEPGroupID),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "local_id", &conn.LocalID),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "peer_ep_group_id", &conn.PeerEPGroupID),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "name", &conn.Name),
					resource.TestCheckResourceAttrPtr("viettelidc_vpnaas_site_connection_v2.conn_1", "dpd.0.action", &conn.DPD.Action),
					resource.TestCheckResourceAttr("viettelidc_vpnaas_site_connection_v2.conn_1", "dpd.0.timeout", strconv.Itoa(conn.DPD.Timeout)),
					resource.TestCheckResourceAttr("viettelidc_vpnaas_site_connection_v2.conn_1", "dpd.0.interval", strconv.Itoa(conn.DPD.Interval)),
				),
			},
		},
	})
}

func testAccCheckSiteConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "viettelidc_vpnaas_site_connection_v2" {
			continue
		}
		_, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Site connection (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckSiteConnectionV2Exists(n string, conn *siteconnections.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		var found *siteconnections.Connection

		found, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*conn = *found

		return nil
	}
}

func testAccSiteConnectionV2Basic() string {
	return fmt.Sprintf(`
	resource "viettelidc_networking_network_v2" "network_1" {
		name           = "tf_test_network"
  		admin_state_up = "true"
	}

	resource "viettelidc_networking_subnet_v2" "subnet_1" {
  		network_id = "${viettelidc_networking_network_v2.network_1.id}"
  		cidr       = "192.168.199.0/24"
  		ip_version = 4
	}

	resource "viettelidc_networking_router_v2" "router_1" {
  		name             = "my_router"
  		external_network_id = "%s"
	}

	resource "viettelidc_networking_router_interface_v2" "router_interface_1" {
  		router_id = "${viettelidc_networking_router_v2.router_1.id}"
  		subnet_id = "${viettelidc_networking_subnet_v2.subnet_1.id}"
	}

	resource "viettelidc_vpnaas_service_v2" "service_1" {
		router_id = "${viettelidc_networking_router_v2.router_1.id}"
		admin_state_up = "false"
	}

	resource "viettelidc_vpnaas_ipsec_policy_v2" "policy_1" {
	}

	resource "viettelidc_vpnaas_ike_policy_v2" "policy_2" {
	}

	resource "viettelidc_vpnaas_endpoint_group_v2" "group_1" {
		type = "cidr"
		endpoints = ["10.0.0.24/24", "10.0.0.25/24"]
	}
	resource "viettelidc_vpnaas_endpoint_group_v2" "group_2" {
		type = "subnet"
		endpoints = [ "${viettelidc_networking_subnet_v2.subnet_1.id}" ]
	}

	resource "viettelidc_vpnaas_site_connection_v2" "conn_1" {
		name = "connection_1"
		ikepolicy_id = "${viettelidc_vpnaas_ike_policy_v2.policy_2.id}"
		ipsecpolicy_id = "${viettelidc_vpnaas_ipsec_policy_v2.policy_1.id}"
		vpnservice_id = "${viettelidc_vpnaas_service_v2.service_1.id}"
		psk = "secret"
		peer_address = "192.168.10.1"
		peer_id = "192.168.10.1"
		local_ep_group_id = "${viettelidc_vpnaas_endpoint_group_v2.group_2.id}"
		peer_ep_group_id = "${viettelidc_vpnaas_endpoint_group_v2.group_1.id}"
		dpd {
			action   = "restart"
			timeout  = 42
			interval = 21
		}
		depends_on = ["viettelidc_networking_router_interface_v2.router_interface_1"]
	}
	`, osExtGwID)
}
