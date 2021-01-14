package search

import (
	_ "net/url"
	_ "strconv"

	pc "github.com/paloaltonetworks/prisma-cloud-go"
)

// ConfigSearch performs a config RQL search.
func ConfigSearch(c pc.PrismaCloudClient, req ConfigRequest) (ConfigResponse, error) {
	c.Log(pc.LogAction, "(get) performing %s", configSingular)

	var resp ConfigResponse

	// Sanity check the time range.
	if err := req.TimeRange.SetType(); err != nil {
		return resp, err
	}

	path := make([]string, 0, len(BaseSuffix)+len(ConfigSuffix))
	path = append(path, BaseSuffix...)
	path = append(path, ConfigSuffix...)

	_, err := c.Communicate("POST", path, nil, req, &resp)
	return resp, err
}

// ConfigSearch performs a config RQL search.
func EventSearch(c pc.PrismaCloudClient, req EventRequest) (EventResponse, error) {
	c.Log(pc.LogAction, "(get) performing %s", eventSingular)

	var resp EventResponse

	// Sanity check the time range.
	if err := req.TimeRange.SetType(); err != nil {
		return resp, err
	}

	path := make([]string, 0, len(BaseSuffix)+len(EventSuffix))
	path = append(path, BaseSuffix...)
	path = append(path, EventSuffix...)

	_, err := c.Communicate("POST", path, nil, req, &resp)
	return resp, err
}

// ConfigSearch performs a config RQL search.
func NetworkSearch(c pc.PrismaCloudClient, req NetworkRequest) (NetworkResponse, error) {
	c.Log(pc.LogAction, "(get) performing %s", networkSingular)

	var resp NetworkResponse

	// Sanity check the time range.
	if err := req.TimeRange.SetType(); err != nil {
		return resp, err
	}

	path := make([]string, 0, len(BaseSuffix)+len(NetworkSuffix))
	path = append(path, BaseSuffix...)
	path = append(path, ConfigSuffix...)

	_, err := c.Communicate("POST", path, nil, req, &resp)
	return resp, err
}