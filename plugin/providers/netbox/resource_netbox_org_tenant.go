package netbox

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/Preskton/go-netbox/netbox/client/tenancy"
	"github.com/Preskton/go-netbox/netbox/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// resourceNetboxOrgTenant is the core Terraform resource structure for the netbox_org_tenant resource.
func resourceNetboxOrgTenant() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxOrgTenantCreate,
		Read:   resourceNetboxOrgTenantRead,
		Update: resourceNetboxOrgTenantUpdate,
		Delete: resourceNetboxOrgTenantDelete,
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
			"tenant_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"comments": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_group_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

// resourceNetboxOrgTenantCreate creates a new Prefix in Netbox.
func resourceNetboxOrgTenantCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	description := d.Get("description").(string)
	comments := d.Get("comments").(string)
	tenantGroupID := int64(d.Get("tenant_group_id").(int))

	var parm = tenancy.NewTenancyTenantsCreateParams().WithData(
		&models.TenantCreateUpdate{
			Name:        &name,
			Slug:        &slug,
			Description: description,
			Comments:    comments,
			Tags:        []string{},
			Group:       tenantGroupID,
			// TODO Tenant Group
		},
	)

	log.Debugf("Executing TenancyTenantsCreate against Netbox: %v", parm)

	out, err := netboxClient.Tenancy.TenancyTenantsCreate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute TenancyTenantsCreate: %v", err)

		return err
	}

	// TODO Probably a better way to parse this ID
	d.SetId(fmt.Sprintf("org/tenant/%d", out.Payload.ID))
	d.Set("tenant_id", out.Payload.ID)

	log.Debugf("Done Executing TenancyTenantsCreate: %v", out)

	return nil
}

// resourceNetboxOrgTenantUpdate applies updates to a Prefix by ID when deltas are detected by Terraform.
func resourceNetboxOrgTenantUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("tenant_id").(int))

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	description := d.Get("description").(string)
	comments := d.Get("comments").(string)
	tenantGroupID := int64(d.Get("tenant_group_id").(int))

	var parm = tenancy.NewTenancyTenantsUpdateParams().
		WithID(id).
		WithData(
			&models.TenantCreateUpdate{
				Name:        &name,
				Slug:        &slug,
				Description: description,
				Comments:    comments,
				Tags:        []string{},
				Group:       tenantGroupID,
			},
		)

	log.Debugf("Executing TenancyTenantsUpdate against Netbox: %v", parm)

	out, err := netboxClient.Tenancy.TenancyTenantsUpdate(parm, nil)

	if err != nil {
		log.Debugf("Failed to execute TenancyTenantsUpdate: %v", err)

		return err
	}

	log.Debugf("Done Executing TenancyTenantsUpdate: %v", out)

	return nil
}

// resourceNetboxOrgTenantRead reads an existing Prefix by ID.
func resourceNetboxOrgTenantRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id := int64(d.Get("tenant_id").(int))

	var readParams = tenancy.NewTenancyTenantsReadParams().WithID(id)

	readResult, err := netboxClient.Tenancy.TenancyTenantsRead(readParams, nil)

	if err != nil {
		log.Debugf("Error fetching Tenant ID # %d from Netbox = %v", id, err)
		return err
	}

	d.Set("name", readResult.Payload.Name)
	d.Set("slug", readResult.Payload.Slug)
	d.Set("description", readResult.Payload.Description)
	d.Set("comments", readResult.Payload.Comments)

	var tenantGroupID int64
	if readResult.Payload.Group != nil {
		tenantGroupID = readResult.Payload.Group.ID
	}
	d.Set("tenant_group_id", tenantGroupID)

	return nil
}

// resourceNetboxOrgTenantDelete deletes an existing Prefix by ID.
func resourceNetboxOrgTenantDelete(d *schema.ResourceData, meta interface{}) error {
	log.Debugf("Deleting Prefix: %v\n", d)

	id := int64(d.Get("tenant_id").(int))

	var deleteParameters = tenancy.NewTenancyTenantsDeleteParams().WithID(id)

	c := meta.(*ProviderNetboxClient).client

	out, err := c.Tenancy.TenancyTenantsDelete(deleteParameters, nil)

	if err != nil {
		log.Debugf("Failed to execute OrgTenantesDelete: %v", err)
	}

	log.Debugf("Done Executing OrgTenantesDelete: %v", out)

	return nil
}
