package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wayfair/terraform-provider-utils/log"
)

const (
	OrganizationEndpointPrefix = "organizations"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

type ForemanOrganization struct {
	// Inherits the base object's attributes
	ForemanObject `json:"foreman_object"`

	Name        string `json:"name"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// ReadOrganization reads the attributes of a ForemanOrganization identified by
// the supplied ID and returns a ForemanOrganization reference.
func (c *Client) ReadOrganization(id int) (*ForemanOrganization, error) {
	log.Tracef("foreman/api/organization.go#Read")

	reqEndpoint := fmt.Sprintf("/%s/%d", OrganizationEndpointPrefix, id)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readOrganization ForemanOrganization
	sendErr := c.SendAndParse(req, &readOrganization)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("readOrganization: [%+v]", readOrganization)

	return &readOrganization, nil
}

// -----------------------------------------------------------------------------
// Query Implementation
// -----------------------------------------------------------------------------

// QueryOrganization queries for a ForemanOrganization based on the attributes
// of the supplied ForemanOrganization reference and returns a QueryResponse
// struct containing query/response metadata and the matching template kinds
func (c *Client) QueryOrganization(t *ForemanOrganization) (QueryResponse, error) {
	log.Tracef("foreman/api/organization.go#Search")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", OrganizationEndpointPrefix)
	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return queryResponse, reqErr
	}

	// dynamically build the query based on the attributes
	reqQuery := req.URL.Query()
	name := `"` + t.Name + `"`
	reqQuery.Set("search", "name="+name)

	req.URL.RawQuery = reqQuery.Encode()
	sendErr := c.SendAndParse(req, &queryResponse)
	if sendErr != nil {
		return queryResponse, sendErr
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	// Results will be Unmarshaled into a []map[string]interface{}
	//
	// Encode back to JSON, then Unmarshal into []ForemanOrganization for
	// the results
	results := []ForemanOrganization{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	jsonDecErr := json.Unmarshal(resultsBytes, &results)
	if jsonDecErr != nil {
		return queryResponse, jsonDecErr
	}
	// convert the search results from []ForemanOrganization to []interface
	// and set the search results on the query
	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}
