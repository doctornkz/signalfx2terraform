package signalfx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/dashboard"
)

// TODO Create simple dashboard

// DashboardAPIURL is the base URL for interacting with dashboard.
const DashboardAPIURL = "/v2/dashboard"

// CreateDashboard creates a dashboard.
func (c *Client) CreateDashboard(dashboardRequest *dashboard.CreateUpdateDashboardRequest) (*dashboard.Dashboard, error) {
	payload, err := json.Marshal(dashboardRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", DashboardAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalDashboard := &dashboard.Dashboard{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboard)

	return finalDashboard, err
}

// DeleteDashboard deletes a dashboard.
func (c *Client) DeleteDashboard(id string) error {
	resp, err := c.doRequest("DELETE", DashboardAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	return nil
}

// GetDashboard gets a dashboard.
func (c *Client) GetDashboard(id string) (*dashboard.Dashboard, error) {
	resp, err := c.doRequest("GET", DashboardAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalDashboard := &dashboard.Dashboard{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboard)

	return finalDashboard, err
}

// UpdateDashboard updates a dashboard.
func (c *Client) UpdateDashboard(id string, dashboardRequest *dashboard.CreateUpdateDashboardRequest) (*dashboard.Dashboard, error) {
	payload, err := json.Marshal(dashboardRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PUT", DashboardAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalDashboard := &dashboard.Dashboard{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboard)

	return finalDashboard, err
}

// SearchDashboard searches for dashboards, given a query string in `name`.
func (c *Client) SearchDashboard(limit int, name string, offset int, tags string) (*dashboard.SearchResult, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest("GET", DashboardAPIURL, params, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalDashboards := &dashboard.SearchResult{}

	err = json.NewDecoder(resp.Body).Decode(finalDashboards)

	return finalDashboards, err
}
