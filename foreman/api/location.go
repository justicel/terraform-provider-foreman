package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wayfair/terraform-provider-utils/log"
)

const (
	LocationEndpointPrefix = "locations"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

type ForemanLocation struct {
	// Inherits the base object's attributes
	ForemanObject
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// ReadLocation reads the attributes of a ForemanLocation identified by
// the supplied ID and returns a ForemanLocation reference.
func (c *Client) ReadLocation(id int) (*ForemanLocation, error) {
	log.Tracef("foreman/api/location.go#Read")

	reqEndpoint := fmt.Sprintf("/%s/%d", LocationEndpointPrefix, id)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readLocation ForemanLocation
	sendErr := c.SendAndParse(req, &readLocation)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("readLocation: [%+v]", readLocation)

	return &readLocation, nil
}

// -----------------------------------------------------------------------------
// Query Implementation
// -----------------------------------------------------------------------------

// QueryLocation queries for a ForemanLocation based on the attributes
// of the supplied ForemanLocation reference and returns a QueryResponse
// struct containing query/response metadata and the matching template kinds
func (c *Client) QueryLocation(t *ForemanLocation) (QueryResponse, error) {
	log.Tracef("foreman/api/location.go#Search")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", LocationEndpointPrefix)
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
	// Encode back to JSON, then Unmarshal into []ForemanLocation for
	// the results
	results := []ForemanLocation{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	jsonDecErr := json.Unmarshal(resultsBytes, &results)
	if jsonDecErr != nil {
		return queryResponse, jsonDecErr
	}
	// convert the search results from []ForemanLocation to []interface
	// and set the search results on the query
	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}
