package nextbus

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

type Client struct {
	URL    string
	client *http.Client
}

func NewClient() (*Client, error) {
	return &Client{
		URL:    "http://webservices.nextbus.com/service/publicXMLFeed",
		client: http.DefaultClient,
	}, nil
}

func (c *Client) AgencyList() (*AgencyListResponse, error) {
	var resp AgencyListResponse
	err := c.command("agencyList", &resp, nil)

	return &resp, err
}

func (c *Client) RouteList(agency string) (*RouteListResponse, error) {
	var resp RouteListResponse
	err := c.command("routeList", &resp, map[string]string{
		"a": agency,
	})

	return &resp, err
}

func (c *Client) RouteConfig(agency, routeTag string) (*RouteConfigResponse, error) {
	var resp RouteConfigResponse
	err := c.command("routeConfig", &resp, map[string]string{
		"a": agency,
		"r": routeTag,
	})

	return &resp, err
}

func (c *Client) PredictionsForStopTag(agency, routeTag, stopTag string) (*PredictionsResponse, error) {
	var resp PredictionsResponse
	err := c.command("predictions", &resp, map[string]string{
		"a": agency,
		"r": routeTag,
		"s": stopTag,
	})

	return &resp, err
}

func (c *Client) PredictionsForStopId(agency, stopId string) (*PredictionsResponse, error) {
	var resp PredictionsResponse
	err := c.command("predictions", &resp, map[string]string{
		"a":      agency,
		"stopId": stopId,
	})

	return &resp, err
}

func (c *Client) Schedule(agency, routeTag string) (*ScheduleResponse, error) {
	return nil, errNotImplemented
}

func (c *Client) Messages(agency string, routeTags ...string) (*MessagesResponse, error) {
	return nil, errNotImplemented
}

func (c *Client) VehicleLocations(agency, routeTag string, lastUpdatedMs uint64) (*VehicleLocationsResponse, error) {
	return nil, errNotImplemented
}

func (c *Client) command(name string, r response, params map[string]string) error {
	req, err := http.NewRequest("GET", c.URL, nil)
	if err != nil {
		return err
	}

	query := url.Values{}
	query.Add("command", name)
	for k, v := range params {
		if v != "" {
			query.Add(k, v)
		}
	}

	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	err = xml.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return err
	}

	err = r.responseError()
	if err != nil {
		return err
	}

	return nil
}
