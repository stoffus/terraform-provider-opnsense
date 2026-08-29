package interfaces_test

import (
	"fmt"
	"testing"

	"github.com/browningluke/terraform-provider-opnsense/internal/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfacesAssignmentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccPreCheck(t) },
		ProtoV6ProviderFactories: acctest.ProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccAssignmentResourceConfig("AssignmentTest"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("opnsense_interfaces_assignment.test", "id"),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "device", "vlan04093"),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "description", "AssignmentTest"),
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "lock", "true"),
					resource.TestCheckResourceAttrPair(
						"data.opnsense_interfaces_assignment.test", "id",
						"opnsense_interfaces_assignment.test", "id",
					),
					resource.TestCheckResourceAttr("data.opnsense_interfaces_assignment.test", "device", "vlan04093"),
				),
			},
			{
				ResourceName:      "opnsense_interfaces_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAssignmentResourceConfig("AssignmentUpdated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("opnsense_interfaces_assignment.test", "description", "AssignmentUpdated"),
				),
			},
		},
	})
}

func testAccAssignmentResourceConfig(description string) string {
	return fmt.Sprintf(`
resource "opnsense_interfaces_vlan" "assignment_test" {
  tag         = 4093
  description = "Assignment test VLAN"
  priority    = 0
  parent      = "vtnet0"
  device      = "vlan04093"
}

resource "opnsense_interfaces_assignment" "test" {
  device      = opnsense_interfaces_vlan.assignment_test.device
  description = %q
}

data "opnsense_interfaces_assignment" "test" {
  id = opnsense_interfaces_assignment.test.id
}
`, description)
}
