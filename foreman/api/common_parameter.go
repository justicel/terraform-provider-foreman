package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wayfair/terraform-provider-utils/log"
)

const (
	CommonParameterEndpointPrefix = "/common_parameters"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

// The ForemanCommonParameter API model represents the commonParameter name. CommonParameters serve as an
// identification string that defines autonomy, authority, or control for
// a portion of a network.
type ForemanCommonParameter struct {
	// Inherits the base object's attributes
	ForemanObject

	// The CommonParameter we actually send
	Name  string `json:"name"`
	Value string `json:"value"`
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// CreateCommonParameter creates a new ForemanCommonParameter with the attributes of the supplied
// ForemanCommonParameter reference and returns the created ForemanCommonParameter reference.
// The returned reference will have its ID and other API default values set by
// this function.
func (c *Client) CreateCommonParameter(d *ForemanCommonParameter) (*ForemanCommonParameter, error) {
	log.Tracef("foreman/api/common_parameter.go#Create")

	reqEndpoint := CommonParameterEndpointPrefix

	// All commonParameters are send individually. Yeay for that
	var createdCommonParameter ForemanCommonParameter
	commonParameterJSONBytes, jsonEncErr := WrapJson("common_parameter", d)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("commonParameterJSONBytes: [%s]", commonParameterJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPost,
		reqEndpoint,
		bytes.NewBuffer(commonParameterJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	sendErr := c.SendAndParse(req, &createdCommonParameter)
	if sendErr != nil {
		return nil, sendErr
	}
	log.Debugf("createdCommonParameter: [%+v]", createdCommonParameter)

	d.Id = createdCommonParameter.Id
	d.Name = createdCommonParameter.Name
	d.Value = createdCommonParameter.Value
	return d, nil
}

// ReadCommonParameter reads the attributes of a ForemanCommonParameter identified by the
// supplied ID and returns a ForemanCommonParameter reference.
func (c *Client) ReadCommonParameter(d *ForemanCommonParameter, id int) (*ForemanCommonParameter, error) {
	log.Tracef("foreman/api/common_parameter.go#Read")

	reqEndpoint := fmt.Sprintf(CommonParameterEndpointPrefix+"/%d", id)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readCommonParameter ForemanCommonParameter
	sendErr := c.SendAndParse(req, &readCommonParameter)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("readCommonParameter: [%+v]", readCommonParameter)

	d.Id = readCommonParameter.Id
	d.Name = readCommonParameter.Name
	d.Value = readCommonParameter.Value
	return d, nil
}

// UpdateCommonParameter deletes all commonParameters for the subject resource and re-creates them
// as we look at them differently on either side this is the safest way to reach sync
func (c *Client) UpdateCommonParameter(d *ForemanCommonParameter, id int) (*ForemanCommonParameter, error) {
	log.Tracef("foreman/api/common_parameter.go#Update")

	reqEndpoint := fmt.Sprintf(CommonParameterEndpointPrefix+"/%d", id)

	commonParameterJSONBytes, jsonEncErr := WrapJson("common_parameter", d)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("commonParameterJSONBytes: [%s]", commonParameterJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPut,
		reqEndpoint,
		bytes.NewBuffer(commonParameterJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var updatedCommonParameter ForemanCommonParameter
	sendErr := c.SendAndParse(req, &updatedCommonParameter)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("updatedCommonParameter: [%+v]", updatedCommonParameter)

	d.Id = updatedCommonParameter.Id
	d.Name = updatedCommonParameter.Name
	d.Value = updatedCommonParameter.Value
	return d, nil
}

// DeleteCommonParameter deletes the ForemanCommonParameters for the given resource
func (c *Client) DeleteCommonParameter(d *ForemanCommonParameter, id int) error {
	log.Tracef("foreman/api/common_parameter.go#Delete")

	reqEndpoint := fmt.Sprintf(CommonParameterEndpointPrefix+"/%d", id)

	req, reqErr := c.NewRequest(
		http.MethodDelete,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return reqErr
	}

	return c.SendAndParse(req, nil)
}

// -----------------------------------------------------------------------------
// Query Implementation
// -----------------------------------------------------------------------------

// QueryCommonParameter queries for a ForemanCommonParameter based on the attributes of the
// supplied ForemanCommonParameter reference and returns a QueryResponse struct
// containing query/response metadata and the matching commonParameters.
func (c *Client) QueryCommonParameter(d *ForemanCommonParameter) (QueryResponse, error) {
	log.Tracef("foreman/api/common_parameter.go#Search")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", CommonParameterEndpointPrefix)
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
	name := `"` + d.Name + `"`
	reqQuery.Set("search", "name="+name)

	req.URL.RawQuery = reqQuery.Encode()
	sendErr := c.SendAndParse(req, &queryResponse)
	if sendErr != nil {
		return queryResponse, sendErr
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	// Results will be Unmarshaled into a []map[string]interface{}
	//
	// Encode back to JSON, then Unmarshal into []ForemanCommonParameter for
	// the results
	results := []ForemanCommonParameter{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	jsonDecErr := json.Unmarshal(resultsBytes, &results)
	if jsonDecErr != nil {
		return queryResponse, jsonDecErr
	}
	// convert the search results from []ForemanCommonParameter to []interface
	// and set the search results on the query
	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}
