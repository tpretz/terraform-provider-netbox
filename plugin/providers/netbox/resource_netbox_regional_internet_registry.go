package netbox

import (
	// "errors"
	// "fmt"
	// "strconv"
	"log"

	"github.com/digitalocean/go-netbox/netbox/client/ipam"
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
		Read:   nil,
		Update: resourceNetboxRegionalInternetRegistryUpdate,
		Delete: resourceNetboxRegionalInternetRegistryDelete,
		Exists: resourceNetboxRegionalInternetRegistryExists,

		Schema: BaseRegionalInternetRegistrySchema(),
	}
}

// Exists is called before Read and obviously makes sure the resource exists.
func resourceNetboxRegionalInternetRegistryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	return true, nil
}

// Create will simply create a new instance of your resource.
// The is also where you will have to set the ID (has to be an Int) of your resource.
// If the API you are using doesn’t provide an ID, you can always use a random Int.
func resourceNetboxRegionalInternetRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Creating RIR: %v\n", d)

	c := meta.(*ProviderNetboxClient).client

	var parm = ipam.NewIPAMRirsCreateParams()

	parm.Data.ID = int64(d.Get("rir_id").(int))
	parm.Data.Name = d.Get("name").(*string)
	parm.Data.Slug = d.Get("slug").(*string)
	parm.Data.IsPrivate = d.Get("is_private").(bool)

	log.Printf("Executing IPAMRirsCreate against Netbox: %v", parm)

	out, err := c.IPAM.IPAMRirsCreate(parm, nil)

	if err != nil {
		log.Printf("Failed to execute IPAMRirsCreate: %v", err)
	}

	log.Printf("Done Executing IPAMRirsCreate: %v", out)

	return nil
}

//Update is optional if your Resource doesn’t support update.
//For example, I’m not using update in the Terraform LDAP Provider.
//I just destroy and recreate the resource everytime there is a change.
func resourceNetboxRegionalInternetRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] JP resourceNetboxRegionalInternetRegistryUpdate: %v\n", d)

	return nil
}

func resourceNetboxRegionalInternetRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Deleting RIR: %v\n", d)

	var deleteParameters = ipam.NewIPAMRirsDeleteParams().
		WithID(int64(d.Get("rir_id").(int)))

	c := meta.(*ProviderNetboxClient).client

	out, err := c.IPAM.IPAMRirsDelete(deleteParameters, nil)

	if err != nil {
		log.Printf("Failed to execute IPAMRirsDelete: %v", err)
	}

	log.Printf("Done Executing IPAMRirsDelete: %v", out)

	return nil
}
