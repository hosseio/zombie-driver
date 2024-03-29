package distance

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

//go:generate moq -out locations_getter_mock.go . LocationsGetter
type LocationsGetter interface {
	GetLocations(driverID string, lastMinutes int) (LocationList, error)
}

type DriverLocationClient struct {
	httpClient *http.Client
	baseUrl    BaseDriverLocationsURL
}

type BaseDriverLocationsURL string

func NewDriverLocationClient(baseURL BaseDriverLocationsURL) DriverLocationClient {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConns:          0,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return DriverLocationClient{client, baseURL}
}

const getDriverLocationsEndpoint = "http://%s/drivers/%s/locations?minutes=%d"

func (c DriverLocationClient) GetLocations(driverID string, lastMinutes int) (LocationList, error) {
	url := fmt.Sprintf(getDriverLocationsEndpoint, c.baseUrl, driverID, lastMinutes)

	var locations LocationList
	if err := c.get(url, &locations); err != nil {
		return nil, err
	}

	return locations, nil
}

func (c *DriverLocationClient) get(url string, payload interface{}) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return ClientErr{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("Invalid status code: %d body: %s", resp.StatusCode, string(body)),
		}
	}

	return json.Unmarshal(body, payload)
}

type ClientErr struct {
	Message    string
	StatusCode int
}

func (e ClientErr) Error() string {
	return e.Message
}
