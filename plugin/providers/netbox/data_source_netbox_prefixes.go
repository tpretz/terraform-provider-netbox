package netbox

import (
	"errors"
	//"fmt"
	"log"
	"strconv"

	// "errors"

	"github.com/Preskton/go-netbox/netbox/client/ipam"
	"github.com/Preskton/go-netbox/netbox/models"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNetboxPrefixes() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceNetboxPrefixesRead,
		Schema: dataSourcePrefixesSchema(),
	}
}

func dataSourceNetboxPrefixParse(d *schema.ResourceData, obj *models.Prefix) {
  d.SetId(strconv.FormatInt(obj.ID, 10))
  d.Set("created", obj.Created.String())
  d.Set("description", obj.Description)
  d.Set("family", obj.Family)
  d.Set("is_pool", obj.IsPool)
  d.Set("prefix", obj.Prefix)
  d.Set("last_updated", obj.LastUpdated)

  if obj.Vlan != nil {
    d.Set("vlan_vid", *obj.Vlan.Vid)
  }

  log.Printf("Finished parsing results from IPAMPrefixesRead")
}

// Read will fetch the data of a resource.
func dataSourceNetboxPrefixesRead(d *schema.ResourceData, meta interface{}) error {
  c := meta.(*ProviderNetboxClient).client

  // primary key lookup, direct
  if id, idOk := d.GetOk("prefixes_id"); idOk {
    parm := ipam.NewIPAMPrefixesReadParams()
		parm.SetID(int64(id.(int)))

		out, err := c.IPAM.IPAMPrefixesRead(parm, nil)

		if err != nil {
      log.Printf("error from IPAMPrefixesRead: %v\n", err)
			return err
    }

    dataSourceNetboxPrefixParse(d, out.Payload)
  } else { // anything else, requires a search
    param := ipam.NewIPAMPrefixesListParams()

    // Add any lookup params

    if vid, vidOk := d.GetOk("vlan_vid"); vidOk {
      vlan_vid := float64(vid.(int))
      param.SetVlanVid(&vlan_vid)
    }

    if query, queryOk := d.GetOk("query"); queryOk {
      query_str := query.(string)
      param.SetQ(&query_str)
    }

    if tenant, tenantOk := d.GetOk("tenant"); tenantOk {
      tenant_str := tenant.(string)
      param.SetTenant(&tenant_str)
    }

    if site, siteOk := d.GetOk("site"); siteOk {
      site_str := site.(string)
      param.SetSite(&site_str)
    }

    if role, roleOk := d.GetOk("role"); roleOk {
      role_str := role.(string)
      param.SetRole(&role_str)
    }

    // limit to 2
    limit := int64(2)
    param.SetLimit(&limit)

		out, err := c.IPAM.IPAMPrefixesList(param, nil)

		if err != nil {
      log.Printf("error from IPAMPrefixesList: %v\n", err)
			return err
    }

    if *out.Payload.Count == 0 {
			return errors.New("Prefix not found")
    } else if *out.Payload.Count > 1 {
      return errors.New("More than one prefix matches search terms, please narrow")
    }

    dataSourceNetboxPrefixParse(d, out.Payload.Results[0])
  }

	return nil
}

func barePrefixesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"prefixes_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"created": &schema.Schema{
			Type: schema.TypeString,
		},
		"description": &schema.Schema{
			Type: schema.TypeString,
		},
		"prefix": &schema.Schema{
			Type: schema.TypeString,
		},
		"family": &schema.Schema{
			Type: schema.TypeString,
		},
		"vlan": &schema.Schema{
			Type: schema.TypeMap,
		},
		"is_pool": &schema.Schema{
			Type: schema.TypeBool,
		},
		"last_updated": &schema.Schema{
			Type: schema.TypeString,
		},
		"vlan_vid": &schema.Schema{
			Type: schema.TypeInt,
		},
		"query": &schema.Schema{
			Type: schema.TypeString,
      Optional: true,
		},
		"tenant": &schema.Schema{
			Type: schema.TypeString,
      Optional: true,
		},
		"site": &schema.Schema{
			Type: schema.TypeString,
      Optional: true,
		},
		"role": &schema.Schema{
			Type: schema.TypeString,
      Optional: true,
		},
	}
}

func resourcePrefixesSchema() map[string]*schema.Schema {
	s := barePrefixesSchema()

	for k, v := range s {
		switch k {
		case "prefixes_id":
			v.Optional = true
			v.ConflictsWith = []string{"vlan_vid"}
		case "prefix":
			v.Optional = true
		case "created":
			v.Optional = true
		case "vlan_vid":
			v.Optional = true
			v.ConflictsWith = []string{"prefixes_id"}
		default:
			v.Computed = true
		}
	}
	// Add the remove_dns_on_delete item to the schema. This is a meta-parameter
	// that is not part of the API resource and exists to instruct NETBOX to
	// gracefully remove the address from its DNS integrations as well when it is
	// removed. The default on this option is true.
	s["remove_dns_on_delete"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  true,
	}
	return s

}

// dataSourceAddressSchema returns the schema for the dataSourceNetboxPrefixes data
// source. It sets the searchable fields and sets up the attribute conflicts
// between IP address and address ID. It also ensures that all fields are
// computed as well.
func dataSourcePrefixesSchema() map[string]*schema.Schema {
	s := barePrefixesSchema()
	for k, v := range s {
		switch k {
		case "prefixes_id":
			v.Optional = true
		case "vlan_vid":
			v.Optional = true
		case "prefix":
			v.Optional = true
		case "created":
			v.Optional = true
			//v.ConflictsWith = []string{"ip_address", "subnet_id", "description", "hostname", "custom_field_filter"}
		default:
			v.Computed = true
		}
	}
	// Add the custom_field_filter item to the schema. This is a meta-parameter
	// that allows searching for a custom field value in the data source.
	s["custom_field_filter"] = customFieldFilterSchema([]string{"prefixes_id"})

	return s
}
