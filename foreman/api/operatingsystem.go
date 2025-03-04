package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wayfair/terraform-provider-utils/log"
)

const (
	OperatingSystemEndpointPrefix = "operatingsystems"
)

// -----------------------------------------------------------------------------
// Struct Definition and Helpers
// -----------------------------------------------------------------------------

// The ForemanOperatingSystem API model represents an operating system
type ForemanOperatingSystem struct {
	// Inherits the base object's attributes
	ForemanObject

	// Title is a computed property by Foreman. The operating system's
	// title is a concatentation of the OS name, major, and minor versions
	// to get a full operating system release.
	Title string `json:"title"`
	// Major release version
	Major string `json:"major"`
	// Minor release version
	Minor string `json:"minor"`
	// Additional information about the operating system
	Description string `json:"description"`
	// Operating sysem family. Available values: AIX, Altlinux, Archlinux,
	// Coreos, Debian, Freebsd, Gentoo, Junos, NXOS, Redhat, Solaris, Suse,
	// Windows.
	Family string `json:"family"`
	// Code name or release name for the specific operating system version
	ReleaseName string `json:"release_name"`
	// Root password hash function to use.  If set, valid values are "MD5",
	// "SHA256", "SHA512", and "Base64"
	PasswordHash string `json:"password_hash"`
	// Provisoning Template Ids
	ProvisioningTemplateIds []int `json:"provisioning_template_ids,omitempty"`
	// Media Ids
	MediumIds []int `json:"medium_ids,omitempty"`
	// Architecture Ids
	ArchitectureIds []int `json:"architecture_ids,omitempty"`
	// Partitiontable Ids
	PartitiontableIds []int `json:"ptable_ids,omitempty"`
}

// ForemanOperating struct used for JSON decode.  Foreman API returns the ids
// back as a list of ForemanObjects with some of the attributes of the data
// types. However, we are only interested in the IDs returned.
type foremanOsRespJSON struct {
	ProvisioningTemplates []ForemanObject `json:"provisioning_templates"`
	Media                 []ForemanObject `json:"media"`
	Architectures         []ForemanObject `json:"architectures"`
	Partitiontables       []ForemanObject `json:"ptables"`
}

// Implement the Unmarshaler interface
func (o *ForemanOperatingSystem) UnmarshalJSON(b []byte) error {
	var jsonDecErr error

	// Unmarshal the common Foreman object properties
	var fo ForemanObject
	jsonDecErr = json.Unmarshal(b, &fo)
	if jsonDecErr != nil {
		return jsonDecErr
	}
	o.ForemanObject = fo
	var foJSON foremanOsRespJSON
	jsonDecErr = json.Unmarshal(b, &foJSON)
	if jsonDecErr != nil {
		return jsonDecErr
	}
	o.ProvisioningTemplateIds = foremanObjectArrayToIdIntArray(foJSON.ProvisioningTemplates)
	o.ArchitectureIds = foremanObjectArrayToIdIntArray(foJSON.Architectures)
	o.MediumIds = foremanObjectArrayToIdIntArray(foJSON.Media)
	o.PartitiontableIds = foremanObjectArrayToIdIntArray(foJSON.Partitiontables)

	var foMap map[string]interface{}
	jsonDecErr = json.Unmarshal(b, &foMap)
	if jsonDecErr != nil {
		return jsonDecErr
	}
	log.Debugf("foMap: [%v]", foMap)

	var ok bool
	if o.Title, ok = foMap["title"].(string); !ok {
		o.Title = ""
	}
	if o.Major, ok = foMap["major"].(string); !ok {
		o.Major = ""
	}
	if o.Minor, ok = foMap["minor"].(string); !ok {
		o.Minor = ""
	}
	if o.Description, ok = foMap["description"].(string); !ok {
		o.Description = ""
	}
	if o.Family, ok = foMap["family"].(string); !ok {
		o.Family = ""
	}
	if o.ReleaseName, ok = foMap["release_name"].(string); !ok {
		o.ReleaseName = ""
	}
	if o.PasswordHash, ok = foMap["password_hash"].(string); !ok {
		o.PasswordHash = ""
	}

	return nil
}

// -----------------------------------------------------------------------------
// CRUD Implementation
// -----------------------------------------------------------------------------

// CreateOperatingSystem creates a new ForemanOperatingSystem with the
// attributes of the supplied ForemanOperatingSystem reference and returns the
// created ForemanOperatingSystem reference.  The returned reference will have
// its ID and other API default values set by this function.
func (c *Client) CreateOperatingSystem(o *ForemanOperatingSystem) (*ForemanOperatingSystem, error) {
	log.Tracef("foreman/api/operatingsystem.go#Create")

	reqEndpoint := fmt.Sprintf("/%s", OperatingSystemEndpointPrefix)

	osJSONBytes, jsonEncErr := WrapJson("operatingsystem", o)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("osJSONBytes: [%s]", osJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPost,
		reqEndpoint,
		bytes.NewBuffer(osJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var createdOperatingSystem ForemanOperatingSystem
	sendErr := c.SendAndParse(req, &createdOperatingSystem)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("createdOperatingSystem: [%+v]", createdOperatingSystem)

	return &createdOperatingSystem, nil
}

// ReadOperatingSystem reads the attributes of a ForemanOperatingSystem
// identified by the supplied ID and returns a ForemanOperatingSystem
// reference.
func (c *Client) ReadOperatingSystem(id int) (*ForemanOperatingSystem, error) {
	log.Tracef("foreman/api/operatingsystem.go#Read")

	reqEndpoint := fmt.Sprintf("/%s/%d", OperatingSystemEndpointPrefix, id)

	req, reqErr := c.NewRequest(
		http.MethodGet,
		reqEndpoint,
		nil,
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var readOperatingSystem ForemanOperatingSystem
	sendErr := c.SendAndParse(req, &readOperatingSystem)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("readOperatingSystem: [%+v]", readOperatingSystem)

	return &readOperatingSystem, nil
}

// UpdateOperatingSystem updates a ForemanOperatingSystem's attributes.  The
// operating system with the ID of the supplied ForemanOperatingSystem will be
// updated. A new ForemanOperatingSystem reference is returned with the
// attributes from the result of the update operation.
func (c *Client) UpdateOperatingSystem(o *ForemanOperatingSystem) (*ForemanOperatingSystem, error) {
	log.Tracef("foreman/api/operatingsystem.go#Update")

	reqEndpoint := fmt.Sprintf("/%s/%d", OperatingSystemEndpointPrefix, o.Id)

	osJSONBytes, jsonEncErr := WrapJson("operatingsystem", o)
	if jsonEncErr != nil {
		return nil, jsonEncErr
	}

	log.Debugf("osJSONBytes: [%s]", osJSONBytes)

	req, reqErr := c.NewRequest(
		http.MethodPut,
		reqEndpoint,
		bytes.NewBuffer(osJSONBytes),
	)
	if reqErr != nil {
		return nil, reqErr
	}

	var updatedOperatingSystem ForemanOperatingSystem
	sendErr := c.SendAndParse(req, &updatedOperatingSystem)
	if sendErr != nil {
		return nil, sendErr
	}

	log.Debugf("updatedOperatingSystem: [%+v]", updatedOperatingSystem)

	return &updatedOperatingSystem, nil
}

// DeleteOperatingSystem deletes the ForemanOperatingSystem identified by the
// supplied ID
func (c *Client) DeleteOperatingSystem(id int) error {
	log.Tracef("foreman/api/operatingsystem.go#Delete")

	reqEndpoint := fmt.Sprintf("/%s/%d", OperatingSystemEndpointPrefix, id)

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

// QueryOperatingSystem queries for a ForemanOperatingSystem based on the
// attributes of the supplied ForemanOperatingSystem reference and returns a
// QueryResponse struct containing query/response metadata and the matching
// operating systems.
func (c *Client) QueryOperatingSystem(o *ForemanOperatingSystem) (QueryResponse, error) {
	log.Tracef("foreman/api/operatingsystem.go#Search")

	queryResponse := QueryResponse{}

	reqEndpoint := fmt.Sprintf("/%s", OperatingSystemEndpointPrefix)

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
	title := `"` + o.Title + `"`
	reqQuery.Set("search", "title="+title)

	req.URL.RawQuery = reqQuery.Encode()
	sendErr := c.SendAndParse(req, &queryResponse)
	if sendErr != nil {
		return queryResponse, sendErr
	}

	log.Debugf("queryResponse: [%+v]", queryResponse)

	// Results will be Unmarshaled into a []map[string]interface{}
	//
	// Encode back to JSON, then Unmarshal into []ForemanOperatingSystem for
	// the results
	results := []ForemanOperatingSystem{}
	resultsBytes, jsonEncErr := json.Marshal(queryResponse.Results)
	if jsonEncErr != nil {
		return queryResponse, jsonEncErr
	}
	jsonDecErr := json.Unmarshal(resultsBytes, &results)
	if jsonDecErr != nil {
		return queryResponse, jsonDecErr
	}
	// conver the search results from []ForemanOperatingSystem to []interface
	// and set the search results on the query
	iArr := make([]interface{}, len(results))
	for idx, val := range results {
		iArr[idx] = val
	}
	queryResponse.Results = iArr

	return queryResponse, nil
}
