package signalfx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/signalfx/signalfx-go/chart"
)

// ChartAPIURL is the base URL for interacting with charts.
const ChartAPIURL = "/v2/chart"

// CreateChart creates a chart.
func (c *Client) CreateChart(chartRequest *chart.CreateUpdateChartRequest) (*chart.Chart, error) {
	payload, err := json.Marshal(chartRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", ChartAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalChart := &chart.Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)

	return finalChart, err
}

// DeleteChart deletes a chart.
func (c *Client) DeleteChart(id string) error {
	resp, err := c.doRequest("DELETE", ChartAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("Unexpected status code: " + resp.Status)
	}

	return nil
}

// GetChart gets a chart.
func (c *Client) GetChart(id string) (*chart.Chart, error) {
	resp, err := c.doRequest("GET", ChartAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalChart := &chart.Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)

	return finalChart, err
}

// UpdateChart updates a chart.
func (c *Client) UpdateChart(id string, chartRequest *chart.CreateUpdateChartRequest) (*chart.Chart, error) {
	payload, err := json.Marshal(chartRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PUT", ChartAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalChart := &chart.Chart{}

	err = json.NewDecoder(resp.Body).Decode(finalChart)

	return finalChart, err
}

// SearchCharts searches for charts, given a query string in `name`.
func (c *Client) SearchCharts(limit int, name string, offset int, tags string) (*chart.SearchResult, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest("GET", ChartAPIURL, params, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalCharts := &chart.SearchResult{}

	err = json.NewDecoder(resp.Body).Decode(finalCharts)

	return finalCharts, err
}
