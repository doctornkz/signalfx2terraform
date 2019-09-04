package signalfx

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// IntegrationAPIURL is the base URL for interacting with intergrations.
const IntegrationAPIURL = "/v2/integration"

// DeleteIntegration deletes an integration.
func (c *Client) DeleteIntegration(id string) error {
	resp, err := c.doRequest("DELETE", IntegrationAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		message, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	return nil
}

// GetIntegration gets a integration.
func (c *Client) GetIntegration(id string) (map[string]interface{}, error) {
	resp, err := c.doRequest("GET", IntegrationAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalIntegration := make(map[string]interface{})

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)

	return finalIntegration, err
}
