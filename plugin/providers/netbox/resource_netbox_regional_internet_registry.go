package netbox

import (
	"log"
	"strconv"

	"github.com/digitalocean/go-netbox/netbox/client/ipam"
	"github.com/digitalocean/go-netbox/netbox/models"
	"github.com/hashicorp/terraform/helper/schema"
)

// resourceNetboxAddress returns the resource structure for the netbox_address
// resource.
//
// Note that we use the data source read function here to pull down data, as
// read workflow is identical for both the resource and the data source.
func resourceNetboxRegionalInternetRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxRegionalInternetRegistryCreate,
		Read:   resourceNetboxRegionalInternetRegistryRead,
		Update: resourceNetboxRegionalInternetRegistryUpdate,
		Delete: resourceNetboxRegionalInternetRegistryDelete,
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
			"is_private": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

// Create will simply create a new instance of your resource.
// The is also where you will have to set the ID (has to be an Int) of your resource.
// If the API you are using doesn’t provide an ID, you can always use a random Int.
func resourceNetboxRegionalInternetRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	isPrivate := d.Get("is_private").(bool)

	var parm = ipam.NewIPAMRirsCreateParams().WithData(
		&models.RIR{
			Slug:      &slug,
			Name:      &name,
			IsPrivate: isPrivate,
		},
	)

	log.Printf("Executing IPAMRirsCreate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMRirsCreate(parm, nil)

	if err != nil {
		log.Printf("Failed to execute IPAMRirsCreate: %v", err)

		return err
	}

	// TODO Probably a better way to parse this ID
	d.SetId(strconv.Itoa(int(out.Payload.ID)))

	log.Printf("Done Executing IPAMRirsCreate: %v", out)

	return nil
}

//Update is optional if your Resource doesn’t support update.
//For example, I’m not using update in the Terraform LDAP Provider.
//I just destroy and recreate the resource everytime there is a change.
func resourceNetboxRegionalInternetRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Printf("Error parsing RIR ID %v = %v", d.Id(), err)
		return err
	}

	name := d.Get("name").(string)
	slug := d.Get("slug").(string)
	isPrivate := d.Get("is_private").(bool)

	var parm = ipam.NewIPAMRirsUpdateParams().
		WithID(int64(id)).
		WithData(
			&models.RIR{
				Slug:      &slug,
				Name:      &name,
				IsPrivate: isPrivate,
			},
		)

	log.Printf("Executing IPAMRirsUpdate against Netbox: %v", parm)

	out, err := netboxClient.IPAM.IPAMRirsUpdate(parm, nil)

	if err != nil {
		log.Printf("Failed to execute IPAMRirsUpdate: %v", err)

		return err
	}

	log.Printf("Done Executing IPAMRirsUpdate: %v", out)

	return nil
}

func resourceNetboxRegionalInternetRegistryRead(d *schema.ResourceData, meta interface{}) error {
	netboxClient := meta.(*ProviderNetboxClient).client

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Printf("Error parsing RIR ID %v = %v", d.Id(), err)
		return err
	}

	var readParams = ipam.NewIPAMRirsReadParams().WithID(int64(id))

	readRirResult, err := netboxClient.IPAM.IPAMRirsRead(readParams, nil)

	if err != nil {
		log.Printf("Error fetching RIR ID # %d from Netbox = %v", id, err)
		return err
	}

	d.Set("name", readRirResult.Payload.Name)
	d.Set("slug", readRirResult.Payload.Slug)
	d.Set("is_private", readRirResult.Payload.IsPrivate)

	return nil
}

func resourceNetboxRegionalInternetRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting RIR: %v\n", d)

	id, err := strconv.Atoi(d.Id())

	if err != nil {
		log.Printf("Error parsing RIR ID %v = %v", d.Id(), err)
		return err
	}

	var deleteParameters = ipam.NewIPAMRirsDeleteParams().WithID(int64(id))

	c := meta.(*ProviderNetboxClient).client

	out, err := c.IPAM.IPAMRirsDelete(deleteParameters, nil)

	if err != nil {
		log.Printf("Failed to execute IPAMRirsDelete: %v", err)
	}

	log.Printf("Done Executing IPAMRirsDelete: %v", out)

	return nil
}
