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

	"github.com/signalfx/signalfx-go/team"
)

// TeamAPIURL is the base URL for interacting with teams.
const TeamAPIURL = "/v2/team"

// CreateTeam creates a team.
func (c *Client) CreateTeam(t *team.CreateUpdateTeamRequest) (*team.Team, error) {
	payload, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("POST", TeamAPIURL, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalTeam := &team.Team{}

	err = json.NewDecoder(resp.Body).Decode(finalTeam)

	return finalTeam, err
}

// DeleteTeam deletes a team.
func (c *Client) DeleteTeam(id string) error {
	resp, err := c.doRequest("DELETE", TeamAPIURL+"/"+id, nil, nil)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("Unexpected status code: " + resp.Status)
	}

	return nil
}

// GetTeam gets a team.
func (c *Client) GetTeam(id string) (*team.Team, error) {
	resp, err := c.doRequest("GET", TeamAPIURL+"/"+id, nil, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalTeam := &team.Team{}

	err = json.NewDecoder(resp.Body).Decode(finalTeam)

	return finalTeam, err
}

// UpdateTeam updates a team.
func (c *Client) UpdateTeam(id string, t *team.CreateUpdateTeamRequest) (*team.Team, error) {
	payload, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	resp, err := c.doRequest("PUT", TeamAPIURL+"/"+id, nil, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		message, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Unexpected status code: %d: %s", resp.StatusCode, message)
	}

	finalTeam := &team.Team{}

	err = json.NewDecoder(resp.Body).Decode(finalTeam)

	return finalTeam, err
}

// SearchTeam searches for teams, given a query string in `name`.
func (c *Client) SearchTeam(limit int, name string, offset int, tags string) (*team.SearchResults, error) {
	params := url.Values{}
	params.Add("limit", strconv.Itoa(limit))
	params.Add("name", name)
	params.Add("offset", strconv.Itoa(offset))
	params.Add("tags", tags)

	resp, err := c.doRequest("GET", TeamAPIURL, params, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	finalTeams := &team.SearchResults{}

	err = json.NewDecoder(resp.Body).Decode(finalTeams)

	return finalTeams, err
}
