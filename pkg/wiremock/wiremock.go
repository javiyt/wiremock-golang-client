//go:generate easyjson $GOFILE

package wiremock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mailru/easyjson"
)

// Client wiremock client instance
type Client struct {
	host   string
	port   uint
	client *http.Client
}

//easyjson:json
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
	} `json:"basicAuthCredentials,omitempty"`
	Cookies      map[string]string `json:"cookies,omitempty"`
	BodyPatterns map[string]string `json:"bodyPatterns,omitempty"`
}

//easyjson:json
type wResponse struct {
	Median                        uint              `json:"median,omitempty"`
	Sigma                         uint              `json:"sigma,omitempty"`
	Type                          string            `json:"type,omitempty"`
	Status                        uint              `json:"status"`
	StatusMessage                 string            `json:"statusMessage,omitempty"`
	Headers                       map[string]string `json:"headers,omitempty"`
	AdditionalProxyRequestHeaders map[string]string `json:"additionalProxyRequestHeaders,omitempty"`
	Body                          string            `json:"body,omitempty"`
	Base64Body                    string            `json:"base64Body,omitempty"`
	JSONBody                      json.RawMessage   `json:"jsonBody,omitempty"`
	BodyFileName                  string            `json:"bodyFileName,omitempty"`
	Fault                         string            `json:"fault,omitempty"`
	FixedDelayMilliseconds        uint              `json:"fixedDelayMilliseconds,omitempty"`
	FromConfiguredStub            bool              `json:"fromConfiguredStub,omitempty"`
	TransformerParameters         map[string]string `json:"transformerParameters,omitempty"`
	Transformers                  []string          `json:"transformers,omitempty"`
}

// Mappings hold mappings configured on wiremock
//easyjson:json
type Mappings struct {
	ID                    string            `json:"id"`
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

// Mapping Main mapping data
//easyjson:json
type Mapping struct {
	Mappings []Mappings `json:"mappings"`
	Meta     struct {
		Total uint `json:"total"`
	} `json:"meta"`
}

// NewWiremockClient generates a new wiremock client instance
func NewWiremockClient(host string, port uint) *Client {
	return &Client{
		host:   host,
		port:   port,
		client: &http.Client{Transport: &http.Transport{}},
	}
}

// Mappings get all mappings defined on wiremock
func (w *Client) Mappings() (Mapping, error) {
	var mapping Mapping

	resp, err := w.client.Get(fmt.Sprintf("http://%s:%v/__admin/mappings", w.host, w.port))
	if err != nil {
		return Mapping{}, err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return Mapping{}, fmt.Errorf("error got from API, status code: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Mapping{}, fmt.Errorf("error reading response body: %w", err)
	}

	easyjson.Unmarshal(body, &mapping)

	return mapping, nil
}
