resource "opnsense_interfaces_vlan" "servers" {
  description = "Servers VLAN"
  tag         = 100
  parent      = "vtnet0"
  device      = "vlan0100"
}

resource "opnsense_interfaces_assignment" "servers" {
  device      = opnsense_interfaces_vlan.servers.device
  description = "Servers"
}
