package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/tpretz/go-netbox/netbox/client/ipam"
	"github.com/tpretz/go-netbox/netbox/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// resourceNetboxIpamVrfDomain is the core Terraform resource structure for the netbox_ipam_vrf_domain resource.
func resourceNetboxIpamVrfDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamVrfDomainCreate,
		Read:   resourceNetboxIpamVrfDomainRead,
		Update: resourceNetboxIpamVrfDomainUpdate,
		Delete: resourceNetboxIpamVrfDomainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"route_distinguisher": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enforce_unique": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"vrf_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tenant_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

// resourceNetboxIpamVrfDomainCreate creates a new VRF in Netbox.
func resourceNetboxIpamVrfDomainCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	name := d.Get("name").(string)
	routeDistinguisher := d.Get("route_distinguisher").(string)
	enforceUnique := d.Get("enforce_unique").(bool)
	description := d.Get("description").(string)
	tenantID := int64(d.Get("tenant_id").(int))

	var parm = ipam.NewIPAMVrfsCreateParams().WithData(
		&models.VRFCreateUpdate{
			Rd:            &routeDistinguisher,
			Name:          &name,
			Description:   description,
			EnforceUnique: enforceUnique,
			Tenant:        tenantID,
			Tags:          []string{},
		},
	)

	log.Debugf("Executing IPAMVrfsCreate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMVrfsCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMVrfsCreate: %v", err)

		return err
	}

	// TODO Probably a better way to parse this ID
	d.SetId(fmt.Sprintf("ipam/vrf/%d", out.Payload.ID))
	d.Set("vrf_id", out.Payload.ID)

	log.Debugf("Done Executing IPAMVrfsCreate: %v", out)

	return nil
}

// resourceNetboxIpamVrfDomainUpdate applies updates to a VRF by ID when deltas are detected by Terraform.
func resourceNetboxIpamVrfDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("vrf_id").(int))

	name := d.Get("name").(string)
	routeDistinguisher := d.Get("route_distinguisher").(string)
	enforceUnique := d.Get("enforce_unique").(bool)
	description := d.Get("description").(string)
	tenantID := int64(d.Get("tenant_id").(int))

	var parm = ipam.NewIPAMVrfsUpdateParams().
		WithID(id).
		WithData(
			&models.VRFCreateUpdate{
				Rd:            &routeDistinguisher,
				Name:          &name,
				Description:   description,
				EnforceUnique: enforceUnique,
				Tenant:        tenantID,
				Tags:          []string{},
			},
		)

	log.Debugf("Executing IPAMVrfsUpdate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMVrfsUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMVrfsUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing IPAMVrfsUpdate: %v", out)

	return nil
}

// resourceNetboxIpamVrfDomainRead reads an existing VRF by ID.
func resourceNetboxIpamVrfDomainRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("vrf_id").(int))

	var readParams = ipam.NewIPAMVrfsReadParams().WithID(id)

	readResult, err := netboxClient.IPAM.IPAMVrfsRead(readParams, nil)

	if err != nil {
		log.Debugf("Error fetching VRF ID # %d from Netbox = %v", id, err)
		return err
	}

	d.Set("name", readResult.Payload.Name)
	d.Set("route_distinguisher", readResult.Payload.Rd)
	d.Set("enforce_unique", readResult.Payload.EnforceUnique)
	d.Set("description", readResult.Payload.Description)

	var tenantID int64
	if readResult.Payload.Tenant != nil {
		tenantID = readResult.Payload.Tenant.ID
	}
	d.Set("tenant_id", tenantID)

	return nil
}

// resourceNetboxIpamVrfDomainDelete deletes an existing VRF by ID.
func resourceNetboxIpamVrfDomainDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting VRF: %v\n", d)

	id := int64(d.Get("vrf_id").(int))

	var deleteParameters = ipam.NewIPAMVrfsDeleteParams().WithID(id)

	c := meta.(*ProviderNetboxClient).client

	out, err := c.IPAM.IPAMVrfsDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMVrfsDelete: %v", err)
	}

	log.Debugf("Done Executing IPAMVrfsDelete: %v", out)

	return nil
}
