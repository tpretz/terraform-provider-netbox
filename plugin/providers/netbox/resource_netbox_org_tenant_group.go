package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Preskton/go-netbox/netbox/client/tenancy"
	"github.com/Preskton/go-netbox/netbox/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// resourceNetboxOrgTenantGroup is the core Terraform resource structure for the netbox_org_tenant_group resource.
func resourceNetboxOrgTenantGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxOrgTenantGroupCreate,
		Read:   resourceNetboxOrgTenantGroupRead,
		Update: resourceNetboxOrgTenantGroupUpdate,
		Delete: resourceNetboxOrgTenantGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"tenant_group_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// resourceNetboxOrgTenantGroupCreate creates a new Prefix in Netbox.
func resourceNetboxOrgTenantGroupCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)

	var parm = tenancy.NewTenancyTenantGroupsCreateParams().WithData(
		&models.TenantGroup{
			Name: &name,
			Slug: &slug,
		},
	)

	log.Debugf("Executing TenancyTenantGroupsCreate against Netbox: %v", parm)

	out, err := netboxClient.Tenancy.TenancyTenantGroupsCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute TenancyTenantGroupsCreate: %v", err)

		return err
	}

	// TODO Probably a better way to parse this ID
	d.SetId(fmt.Sprintf("org/tenant-group/%d", out.Payload.ID))
	d.Set("tenant_group_id", out.Payload.ID)

	log.Debugf("Done Executing TenancyTenantGroupsCreate: %v", out)

	return nil
}

// resourceNetboxOrgTenantGroupUpdate applies updates to a Prefix by ID when deltas are detected by Terraform.
func resourceNetboxOrgTenantGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("tenant_group_id").(int))

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)

	var parm = tenancy.NewTenancyTenantGroupsUpdateParams().
		WithID(id).
		WithData(
			&models.TenantGroup{
				Name: &name,
				Slug: &slug,
			},
		)

	log.Debugf("Executing TenancyTenantGroupsUpdate against Netbox: %v", parm)

	out, err := netboxClient.Tenancy.TenancyTenantGroupsUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute TenancyTenantGroupsUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing TenancyTenantGroupsUpdate: %v", out)

	return nil
}

// resourceNetboxOrgTenantGroupRead reads an existing Prefix by ID.
func resourceNetboxOrgTenantGroupRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("tenant_group_id").(int))

	var readParams = tenancy.NewTenancyTenantGroupsReadParams().WithID(id)

	readResult, err := netboxClient.Tenancy.TenancyTenantGroupsRead(readParams, nil)

	if err != nil {
		log.Debugf("Error fetching TenantGroup ID # %d from Netbox = %v", id, err)
		return err
	}

	d.Set("name", readResult.Payload.Name)
	d.Set("slug", readResult.Payload.Slug)

	return nil
}

// resourceNetboxOrgTenantGroupDelete deletes an existing Prefix by ID.
func resourceNetboxOrgTenantGroupDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting TenantGroup: %v\n", d)

	id := int64(d.Get("tenant_group_id").(int))

	var deleteParameters = tenancy.NewTenancyTenantGroupsDeleteParams().WithID(id)

	c := meta.(*ProviderNetboxClient).client

	out, err := c.Tenancy.TenancyTenantGroupsDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute OrgTenantGroupsDelete: %v", err)
	}

	log.Debugf("Done Executing OrgTenantGroupsDelete: %v", out)

	return nil
}
