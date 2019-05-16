package foreman

import (
	"fmt"
	"strconv"

	"github.com/HanseMerkur/terraform-provider-foreman/foreman/api"
	"github.com/wayfair/terraform-provider-utils/autodoc"
	"github.com/wayfair/terraform-provider-utils/log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceForemanLocation() *schema.Resource {
	return &schema.Resource{

		Read: dataSourceForemanLocationRead,

		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				Description: fmt.Sprintf(
					"Name of location. %s",
					autodoc.MetaExample,
				),
			},
		},
	}
}

// -----------------------------------------------------------------------------
// Conversion Helpers
// -----------------------------------------------------------------------------

// buildForemanLocation constructs a ForemanLocation reference from a
// resource data reference.  The struct's  members are populated from the data
// populated in the resource data.  Missing members will be left to the zero
// value for that member's type.
func buildForemanLocation(d *schema.ResourceData) *api.ForemanLocation {
	t := api.ForemanLocation{}
	obj := buildForemanObject(d)
	t.ForemanObject = *obj
	return &t
}

// setResourceDataFromForemanLocation sets a ResourceData's attributes from
// the attributes of the supplied ForemanLocation reference
func setResourceDataFromForemanLocation(d *schema.ResourceData, fk *api.ForemanLocation) {
	d.SetId(strconv.Itoa(fk.Id))
	d.Set("name", fk.Name)
}

// -----------------------------------------------------------------------------
// Resource CRUD Operations
// -----------------------------------------------------------------------------

func dataSourceForemanLocationRead(d *schema.ResourceData, meta interface{}) error {
	log.Tracef("data_source_foreman_location.go#Read")

	client := meta.(*api.Client)
	t := buildForemanLocation(d)

	log.Debugf("ForemanLocation: [%+v]", t)

	queryResponse, queryErr := client.QueryLocation(t)
	if queryErr != nil {
		return queryErr
	}

	if queryResponse.Subtotal == 0 {
		return fmt.Errorf("Data source location returned no results")
	} else if queryResponse.Subtotal > 1 {
		return fmt.Errorf("Data source location returned more than 1 result")
	}

	var queryLocation api.ForemanLocation
	var ok bool
	if queryLocation, ok = queryResponse.Results[0].(api.ForemanLocation); !ok {
		return fmt.Errorf(
			"Data source results contain unexpected type. Expected "+
				"[api.ForemanLocation], got [%T]",
			queryResponse.Results[0],
		)
	}
	t = &queryLocation

	log.Debugf("ForemanLocation: [%+v]", t)

	setResourceDataFromForemanLocation(d, t)

	return nil
}
