package netbox

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func BaseRegionalInternetRegistrySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"rir_id": &schema.Schema{
			Type:     schema.TypeInt,
			Optional: false,
		},
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Optional: false,
		},
		"slug": &schema.Schema{
			Type:     schema.TypeString,
			Optional: false,
		},
	}
}
