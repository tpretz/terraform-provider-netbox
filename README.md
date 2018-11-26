[![Build Status](https://travis-ci.com/Preskton/terraform-provider-netbox.svg?branch=master)](https://travis-ci.com/Preskton/terraform-provider-netbox) [![GitHub license](https://img.shields.io/github/license/Preskton/terraform-provider-netbox.svg)](https://github.com/Preskton/terraform-provider-netbox/blob/master/LICENSE) [![GitHub release](https://img.shields.io/github/release/Preskton/terraform-provider-netbox.svg)](https://github.com/Preskton/terraform-provider-netbox/releases/)

# Terraform Provider Plugin for Netbox

This repository holds a external plugin for a [Terraform][1] provider to manage resources within DigitalOcean's [Netbox][2] by way of the the golang API client for Netbox, [go-netbox][3].

[1]: https://www.terraform.io/
[2]: https://github.com/digitalocean/netbox
[3]: https://github.com/Preskton/go-netbox

## About Netbox

[Netbox][2] is an IP address management (IPAM) and data center infrastructure management (DCIM) created by DigitalOcean. By leveraging the work at [go-netbox][3], 
`terraform-provider-netbox` allows you to declaratively describe your infrastructure using HCL to keep track of your infrastructure. The real value of this
solution comes through when you combine it with your other Terraform providers to store information like cloud provider-assigned networks and IPs.

## Installing

See the [Plugin Basics][5] page of the Terraform docs to see how to install this plugin. Check the [releases page][6] to download binaries for
Linux, OS X, and Windows. You'll need to remove the OS & processor architecture from the file name for Terraform to recognize the plugin. Ex: if you are using
the `linux-amd64-terraform-provider-netbox`, you'd rename the file to `terraform-provider-netbox`.

[5]: https://www.terraform.io/docs/plugins/basics.html
[6]: https://github.com/Preskton/terraform-provider-netbox/releases

## Usage

Add the `netbox` provider to your `tf` file like so:

```hcl
provider "netbox" {
    app_id = "abcdef12345678900987654321fedcba"
    endpoint = "https://netbox.tonikensa.splatnet"
}
```

Where `app_id` is a Netbox token created in the Netbox Admin portal (click your username in the top right -> Admin -> Tokens) and `endpoint` is a URI to your Netbox instance (do not include `/api`).

Once configured, you can use any of the following resources:

- IPAM Resources:
  - `netbox_ipam_rir` - regional internet registries
  - `netbox_ipam_vrf` - virtual routing & forwarding groups
  - `netbox_ipam_aggregate` - top level aggregates
  - `netbox_ipam_prefix` - subnet prefixes
  - `netbox_ipam_ip_address` - specific IP addresses
- Organization Resources:
  - `netbox_org_tenant_group` - tenant groups
  - `netbox_org_tenant` - tenants

## Annotated Example

The following is an example that exercises the currently available functionality:

```hcl
provider "netbox" {
    app_id = "abcdef12345678900987654321fedcba"
    endpoint = "https://netbox.tonikensa.splatnet"
}

// Creates a tenant group we can place our tenants in
resource "netbox_org_tenant_group" "splatoon" {
    name = "Splatoon Tenants"
    slug = "splatoon"
}

// Creates a tenant that we can later assign things like circuits, racks, and IPs to (once we build those providers, ha)
resource "netbox_org_tenant" "squid-kids" {
    name = "Squid Kids"
    slug = "squid-kids"
    description = "Squid kids only."
    comments = "This tenant reserved for squid kids only. Should NOT be used for Octolings."
    // Use the tenant group we just made
    tenant_group_id = "${netbox_org_tenant_group.splatoon.tenant_group_id}"
}

// Creates a regional internet registry that is responsible for managing the various addresses we'll be registering
resource "netbox_ipam_rir" "squidland" {
    name = "Squidland IP Addressing Protectorate"
    slug = "squidland"
    is_private = "true"
}

// Creates a Virtual Routing & Forwarding domain
resource "netbox_ipam_vrf" "toni-kensa-west" {
    name = "Toni Kensa GmbH Private Networks"
    route_distinguisher = "toni-kensa-west"
    // Forces all prefixes and IPs to be non-overlapping and unique
    enforce_unique = true
}

// Creates a top level aggregate in which underlying prefixes and IPs will live
resource "netbox_ipam_aggregate" "splatnet" {
    prefix = "192.168.0.0/16"
    description = "Squidland Splatnet"
    // Use the RIR we created earlier
    rir_id = "${netbox_ipam_rir.squidland.rir_id}"
}

// Creates a subnet prefix
resource "netbox_ipam_prefix" "toni-kensa-west-primary" {
    prefix = "192.168.100.0/24"
    description = "Toni Kensa West - Primary Network"
    // Use the VRF we just created
    vrf_id = "${netbox_ipam_vrf.toni-kensa-west.vrf_id}"
    is_pool = true    
}

// Creates another subnet prefix
resource "netbox_ipam_prefix" "toni-kensa-west-secondary" {
    prefix = "192.168.101.0/24"
    description = "Toni Kensa West - Secondary Network"
    vrf_id = "${netbox_ipam_vrf.toni-kensa-west.vrf_id}"
    is_pool = true    
}

// Creates an internal IP address that is "active" (status 1)
resource "netbox_ipam_ip_address" "toni-kensa-west-primary-router" {
    // Still use full CIDR notation for IPs!
    address = "192.168.100.1/32"
    description = "Toni Kensa West Primary Router"
    // Use the VRF from above
    vrf_id = "${netbox_ipam_vrf.toni-kensa-west.vrf_id}"
    // Use the tenant from above
    tenant_id = "${netbox_org_tenant.squid-kids.tenant_id}"
    // Sorry not quite using the names yet - you've got to reference the ID!
    // for IPAM resources, these are available at https://your-netbox/api/ipam/_choices
    status = 1
}

// Creates an "outside" IP address that we NAT our previous internal IP through
resource "netbox_ipam_ip_address" "toni-kensa-west-primary-external" {
    address = "3.3.3.3/32"
    description = "Toni Kensa West Primary External IP"
    vrf_id = "${netbox_ipam_vrf.toni-kensa-west.vrf_id}"
    tenant_id = "${netbox_org_tenant.squid-kids.tenant_id}"
    status = 1
    // This IP (3.3.3.3) NATs for the IP specified here (192.168.100.1)
    nat_inside_ip_address_id = "${netbox_ipam_ip_address.toni-kensa-west-primary-router.ip_address_id}"
}
```

## Copyright Notice

```
Copyright 2018 BB, Inc.
Portions copyright 2018 Preston Doster.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
