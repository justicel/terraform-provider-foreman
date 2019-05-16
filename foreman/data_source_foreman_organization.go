package foreman

import (
	"fmt"
	"strconv"

	"github.com/HanseMerkur/terraform-provider-foreman/foreman/api"
	"github.com/wayfair/terraform-provider-utils/autodoc"
	"github.com/wayfair/terraform-provider-utils/log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceForemanOrganization() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceForemanOrganizationRead,

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Name of organization. %s",
					autodoc.MetaExample,
				),
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Conversion Helpers
// -----------------------------------------------------------------------------

// buildForemanOrganization constructs a ForemanOrganization reference from a
// resource data reference.  The struct's  members are populated from the data
// populated in the resource data.  Missing members will be left to the zero
// value for that member's type.
func buildForemanOrganization(d *schema.ResourceData) *api.ForemanOrganization {
	t := api.ForemanOrganization{}
	obj := buildForemanObject(d)
	t.ForemanObject = *obj
	return &t
}

// setResourceDataFromForemanOrganization sets a ResourceData's attributes from
// the attributes of the supplied ForemanOrganization reference
func setResourceDataFromForemanOrganization(d *schema.ResourceData, fk *api.ForemanOrganization) {
	d.SetId(strconv.Itoa(fk.Id))
	d.Set("name", fk.Name)
}

// -----------------------------------------------------------------------------
// Resource CRUD Operations
// -----------------------------------------------------------------------------

func dataSourceForemanOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("data_source_foreman_location.go#Read")

	client := meta.(*api.Client)
	t := buildForemanOrganization(d)

	log.Debugf("ForemanOrganization: [%+v]", t)

	queryResponse, queryErr := client.QueryOrganization(t)
	if queryErr != nil {
		return queryErr
	}

	if queryResponse.Subtotal == 0 {
		return fmt.Errorf("Data source organization returned no results")
	} else if queryResponse.Subtotal > 1 {
		return fmt.Errorf("Data source organization returned more than 1 result")
	}

	var queryOrganization api.ForemanOrganization
	var ok bool
	if queryOrganization, ok = queryResponse.Results[0].(api.ForemanOrganization); !ok {
		return fmt.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanOrganization], got [%T]",
			queryResponse.Results[0],
		)
	}
	t = &queryOrganization

	log.Debugf("ForemanOrganization: [%+v]", t)

	setResourceDataFromForemanOrganization(d, t)

	return nil
}
