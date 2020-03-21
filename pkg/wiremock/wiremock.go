package wiremock

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	host   string
	port   uint
	client *http.Client
}

type wRequest struct {
	Method               string            `json:"method"`
	URL                  string            `json:"url"`
	URLPath              string            `json:"urlPath,omitempty"`
	URLPathPattern       string            `json:"urlPathPattern,omitempty"`
	URLPattern           string            `json:"urlPattern,omitempty"`
	QueryParameters      map[string]string `json:"queryParameters,omitempty"`
	Headers              map[string]string `json:"headers,omitempty"`
	BasicAuthCredentials struct {
		Password string `json:"password"`
		Username string `json:"username"`
	} `json:"queryParameters,omitempty"`
	Cookies      map[string]string `json:"cookies,omitempty"`
	BodyPatterns map[string]string `json:"bodyPatterns,omitempty"`
}

type wResponse struct {
	Median       uint   `json:"median,omitempty"`
	Sigma        uint   `json:"sigma,omitempty"`
	Status       uint   `json:"status"`
	Body         string `json:"body,omitempty"`
	BodyFileName string `json:"bodyFileName,omitempty"`
}

type Mappings struct {
	Id                    string            `json:"id"`
	UUID                  string            `json:"uuid,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Request               wRequest          `json:"request"`
	Response              wResponse         `json:"response"`
	Persistent            bool              `json:"persistent,omitempty"`
	Priority              uint              `json:"priority,omitempty"`
	ScenarioName          string            `json:"scenarioName,omitempty"`
	RequiredScenarioState string            `json:"requiredScenarioState,omitempty"`
	NewScenarioState      string            `json:"newScenarioState,omitempty"`
	PostServeActions      map[string]string `json:"postServeActions,omitempty"`
	Metadata              map[string]string `json:"metadata,omitempty"`
}

type Mapping struct {
	Mappings []Mappings `json:"mappings"`
	Meta     struct {
		Total uint `json:"total"`
	} `json:"meta"`
}

func NewWiremockClient(host string, port uint) *Client {
	return &Client{
		host:   host,
		port:   port,
		client: &http.Client{Transport: &http.Transport{}},
	}
}

func (w *Client) Mappings() ([]Mappings, error) {
	var mappings []Mappings
	url := fmt.Sprintf("http://%s:%v/__admin/mappings", w.host, w.port)
	resp, err := w.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error got from API, status code: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	fmt.Printf("%s", body)

	return mappings, nil

}
