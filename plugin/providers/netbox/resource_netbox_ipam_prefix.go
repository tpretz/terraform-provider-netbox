package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Preskton/go-netbox/netbox/client/ipam"
	"github.com/Preskton/go-netbox/netbox/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// resourceNetboxIpamPrefix is the core Terraform resource structure for the netbox_ipam_Prefix_domain resource.
func resourceNetboxIpamPrefix() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxIpamPrefixCreate,
		Read:   resourceNetboxIpamPrefixRead,
		Update: resourceNetboxIpamPrefixUpdate,
		Delete: resourceNetboxIpamPrefixDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"prefix": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"prefix_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vrf_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"is_pool": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
		},
	}
}

// resourceNetboxIpamPrefixCreate creates a new Prefix in Netbox.
func resourceNetboxIpamPrefixCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	prefix := d.Get("prefix").(string)
	description := d.Get("description").(string)
	//vrfID := int64(d.Get("vrf_id").(int))
	isPool := d.Get("is_pool").(bool)
	//status := d.Get("status").(string)

	var parm = ipam.NewIPAMPrefixesCreateParams().WithData(
		&models.Prefix{
			Prefix:      &prefix,
			Description: description,
			IsPool:      isPool,
			Tags:        []string{},
		},
	)

	log.Debugf("Executing IPAMPrefixesCreate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMPrefixesCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMPrefixesCreate: %v", err)

		return err
	}

	// TODO Probably a better way to parse this ID
	d.SetId(fmt.Sprintf("ipam/prefix/%d", out.Payload.ID))
	d.Set("prefix_id", out.Payload.ID)

	log.Debugf("Done Executing IPAMPrefixesCreate: %v", out)

	return nil
}

// resourceNetboxIpamPrefixUpdate applies updates to a Prefix by ID when deltas are detected by Terraform.
func resourceNetboxIpamPrefixUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("prefix_id").(int))

	prefix := d.Get("prefix").(string)
	description := d.Get("description").(string)
	//vrfID := int64(d.Get("vrf_id").(int))
	isPool := d.Get("is_pool").(bool)
	//status := d.Get("status").(string)

	var parm = ipam.NewIPAMPrefixesUpdateParams().
		WithID(id).
		WithData(
			&models.Prefix{
				Prefix:      &prefix,
				Description: description,
				IsPool:      isPool,
				Tags:        []string{},
			},
		)

	log.Debugf("Executing IPAMPrefixesUpdate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMPrefixesUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMPrefixesUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing IPAMPrefixesUpdate: %v", out)

	return nil
}

// resourceNetboxIpamPrefixRead reads an existing Prefix by ID.
func resourceNetboxIpamPrefixRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("prefix_id").(int))

	var readParams = ipam.NewIPAMPrefixesReadParams().WithID(id)

	readResult, err := netboxClient.IPAM.IPAMPrefixesRead(readParams, nil)

	if err != nil {
		log.Debugf("Error fetching Prefix ID # %d from Netbox = %v", id, err)
		return err
	}

	var vrfID int64
	if readResult.Payload.Vrf != nil {
		vrfID = readResult.Payload.Vrf.ID
	}

	d.Set("prefix", readResult.Payload.Prefix)
	d.Set("description", readResult.Payload.Description)
	d.Set("vrf_id", vrfID)
	d.Set("is_pool", readResult.Payload.IsPool)

	return nil
}

// resourceNetboxIpamPrefixDelete deletes an existing Prefix by ID.
func resourceNetboxIpamPrefixDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting Prefix: %v\n", d)

	id := int64(d.Get("prefix_id").(int))

	var deleteParameters = ipam.NewIPAMPrefixesDeleteParams().WithID(id)

	c := meta.(*ProviderNetboxClient).client

	out, err := c.IPAM.IPAMPrefixesDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute IPAMPrefixesDelete: %v", err)
	}

	log.Debugf("Done Executing IPAMPrefixesDelete: %v", out)

	return nil
}
