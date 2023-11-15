package smartevse

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func FromHTTP(smartEvseUrl string) (*SmartEVSESettings, error) {
	u, err := url.Parse(smartEvseUrl)
	if err != nil {
		return nil, fmt.Errorf("could not parse URL: %v", err)
	}

	u.Path = "/settings"
	response, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("could GET from URL: %v", err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	settings, err := FromJSON(body)
	if err != nil {
		return nil, err
	}

	return settings, nil
}
