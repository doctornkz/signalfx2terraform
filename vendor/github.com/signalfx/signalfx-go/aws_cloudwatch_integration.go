package signalfx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/signalfx/signalfx-go/integration"
)

// CreateAWSCloudWatchIntegration creates an AWS CloudWatch integration.
func (c *Client) CreateAWSCloudWatchIntegration(acwi *integration.AwsCloudWatchIntegration) (*integration.AwsCloudWatchIntegration, error) {
	payload, err := json.Marshal(acwi)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", IntegrationAPIURL, nil, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalIntegration := integration.AwsCloudWatchIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)

	return &finalIntegration, err
}

// GetAWSCloudWatchIntegration retrieves an AWS CloudWatch integration.
func (c *Client) GetAWSCloudWatchIntegration(id string) (*integration.AwsCloudWatchIntegration, error) {
	resp, err := c.doRequest("GET", IntegrationAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalIntegration := integration.AwsCloudWatchIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)

	return &finalIntegration, err
}

// UpdateAWSCloudWatchIntegration updates an AWS CloudWatch integration.
func (c *Client) UpdateAWSCloudWatchIntegration(id string, acwi *integration.AwsCloudWatchIntegration) (*integration.AwsCloudWatchIntegration, error) {
	payload, err := json.Marshal(acwi)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PUT", IntegrationAPIURL+"/"+id, nil, bytes.NewReader(payload))

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalIntegration := integration.AwsCloudWatchIntegration{}

	err = json.NewDecoder(resp.Body).Decode(&finalIntegration)

	return &finalIntegration, err
}

// DeleteAWSCloudWatchIntegration deletes an AWS CloudWatch integration.
func (c *Client) DeleteAWSCloudWatchIntegration(id string) error {
	resp, err := c.doRequest("DELETE", IntegrationAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		message, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	return err
}
