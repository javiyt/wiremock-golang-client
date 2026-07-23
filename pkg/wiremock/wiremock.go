package wiremock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/json-iterator/go"
)

// Client WireMock client instance
type Client struct {
	host   string
	port   uint
	client *http.Client
}

// Request holds the request matcher data for a WireMock mapping.
type Request struct {
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

// Response holds the response data for a WireMock mapping.
type Response struct {
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

// Mappings hold mappings configured on WireMock
type Mappings struct {
	ID                    string            `json:"id"`
	UUID                  string            `json:"uuid,omitempty"`
	Name                  string            `json:"name,omitempty"`
	Request               Request           `json:"request"`
	Response              Response          `json:"response"`
	Persistent            bool              `json:"persistent,omitempty"`
	Priority              uint              `json:"priority,omitempty"`
	ScenarioName          string            `json:"scenarioName,omitempty"`
	RequiredScenarioState string            `json:"requiredScenarioState,omitempty"`
	NewScenarioState      string            `json:"newScenarioState,omitempty"`
	PostServeActions      map[string]string `json:"postServeActions,omitempty"`
	Metadata              map[string]string `json:"metadata,omitempty"`
}

// Mapping Main mapping data
type Mapping struct {
	Mappings []Mappings `json:"mappings"`
	Meta     struct {
		Total uint `json:"total"`
	} `json:"meta"`
}

// MappingsOptions holds query parameters for listing mappings.
type MappingsOptions struct {
	Limit  *uint
	Offset *uint
}

// MetadataFilter holds a WireMock metadata matcher request body.
type MetadataFilter map[string]any

// ImportMappingsRequest holds the mappings import request body.
type ImportMappingsRequest struct {
	Mappings      []Mappings             `json:"mappings"`
	ImportOptions *ImportMappingsOptions `json:"importOptions,omitempty"`
}

// ImportMappingsOptions holds WireMock mappings import options.
type ImportMappingsOptions struct {
	DuplicatePolicy      string `json:"duplicatePolicy,omitempty"`
	DeleteAllNotInImport bool   `json:"deleteAllNotInImport,omitempty"`
}

const (
	// DuplicatePolicyOverwrite overwrites an existing mapping with the same ID during import.
	DuplicatePolicyOverwrite = "OVERWRITE"
	// DuplicatePolicyIgnore leaves an existing mapping with the same ID unchanged during import.
	DuplicatePolicyIgnore = "IGNORE"
)

const (
	mappingsPath                 = "/__admin/mappings"
	mappingsResetPath            = mappingsPath + "/reset"
	mappingsSavePath             = mappingsPath + "/save"
	mappingsImportPath           = mappingsPath + "/import"
	mappingsFindByMetadataPath   = mappingsPath + "/find-by-metadata"
	mappingsRemoveByMetadataPath = mappingsPath + "/remove-by-metadata"
	mappingsUnmatchedPath        = mappingsPath + "/unmatched"
)

// NewWireMockClient generates a new WireMock client instance
func NewWireMockClient(host string, port uint, client *http.Client) *Client {
	if client == nil {
		client = &http.Client{Transport: &http.Transport{}}
	}

	return &Client{
		host:   host,
		port:   port,
		client: client,
	}
}

// Mappings get all mappings defined on WireMock
func (w *Client) Mappings() (Mapping, error) {
	return w.MappingsWithOptions(MappingsOptions{})
}

// MappingsWithOptions get mappings defined on WireMock with optional pagination.
func (w *Client) MappingsWithOptions(options MappingsOptions) (Mapping, error) {
	var mapping Mapping

	path := mappingsPath
	query := url.Values{}
	if options.Limit != nil {
		query.Set("limit", strconv.FormatUint(uint64(*options.Limit), 10))
	}
	if options.Offset != nil {
		query.Set("offset", strconv.FormatUint(uint64(*options.Offset), 10))
	}
	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	err := w.doJSON(http.MethodGet, path, nil, http.StatusOK, &mapping)
	return mapping, err
}

// SaveMapping stores a new mapping in WireMock.
func (w *Client) SaveMapping(mapping Mappings) (Mappings, error) {
	var savedMapping Mappings

	err := w.doJSON(http.MethodPost, mappingsPath, mapping, http.StatusCreated, &savedMapping)
	return savedMapping, err
}

// DeleteMappings deletes all stub mappings.
func (w *Client) DeleteMappings() error {
	return w.do(http.MethodDelete, mappingsPath, nil, http.StatusOK)
}

// ResetMappings restores stub mappings to the defaults defined in the backing store.
func (w *Client) ResetMappings() error {
	return w.do(http.MethodPost, mappingsResetPath, nil, http.StatusOK)
}

// PersistMappings saves all persistent stub mappings to the backing store.
func (w *Client) PersistMappings() error {
	return w.do(http.MethodPost, mappingsSavePath, nil, http.StatusOK)
}

// ImportMappings imports stub mappings to WireMock.
func (w *Client) ImportMappings(request ImportMappingsRequest) error {
	return w.doJSON(http.MethodPost, mappingsImportPath, request, http.StatusOK, nil)
}

// Mapping gets a stub mapping by ID.
func (w *Client) Mapping(stubMappingID string) (Mappings, error) {
	var mapping Mappings

	err := w.doJSON(http.MethodGet, mappingPath(stubMappingID), nil, http.StatusOK, &mapping)
	return mapping, err
}

// UpdateMapping updates a stub mapping by ID.
func (w *Client) UpdateMapping(stubMappingID string, mapping Mappings) (Mappings, error) {
	var updatedMapping Mappings

	err := w.doJSON(http.MethodPut, mappingPath(stubMappingID), mapping, http.StatusOK, &updatedMapping)
	return updatedMapping, err
}

// DeleteMapping deletes a stub mapping by ID.
func (w *Client) DeleteMapping(stubMappingID string) error {
	return w.do(http.MethodDelete, mappingPath(stubMappingID), nil, http.StatusOK)
}

// FindMappingsByMetadata finds stub mappings by matching on their metadata.
func (w *Client) FindMappingsByMetadata(filter MetadataFilter) (Mapping, error) {
	var mapping Mapping

	err := w.doJSON(http.MethodPost, mappingsFindByMetadataPath, filter, http.StatusOK, &mapping)
	return mapping, err
}

// RemoveMappingsByMetadata deletes stub mappings matching metadata.
func (w *Client) RemoveMappingsByMetadata(filter MetadataFilter) error {
	return w.doJSON(http.MethodPost, mappingsRemoveByMetadataPath, filter, http.StatusOK, nil)
}

// UnmatchedMappings gets stub mappings that have not matched any requests in the journal.
func (w *Client) UnmatchedMappings() (Mapping, error) {
	var mapping Mapping

	err := w.doJSON(http.MethodGet, mappingsUnmatchedPath, nil, http.StatusOK, &mapping)
	return mapping, err
}

// DeleteUnmatchedMappings deletes stub mappings that have not matched any requests in the journal.
func (w *Client) DeleteUnmatchedMappings() error {
	return w.do(http.MethodDelete, mappingsUnmatchedPath, nil, http.StatusOK)
}

func (w *Client) doJSON(method, path string, requestBody any, expectedStatusCode int, responseBody any) error {
	var body io.Reader
	if requestBody != nil {
		data, err := jsoniter.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("error marshaling request body %w", err)
		}
		body = bytes.NewReader(data)
	}

	err := w.do(method, path, body, expectedStatusCode, responseBody)
	if err != nil {
		return err
	}

	return nil
}

func (w *Client) do(method, path string, body io.Reader, expectedStatusCode int, responseBody ...any) error {
	req, err := http.NewRequest(method, w.url(path), body)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != expectedStatusCode {
		return fmt.Errorf("error got from API, status code: %v", resp.StatusCode)
	}

	if len(responseBody) == 0 || responseBody[0] == nil {
		return nil
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	err = jsoniter.Unmarshal(responseData, responseBody[0])
	if err != nil {
		return fmt.Errorf("error unmarshaling response %w", err)
	}

	return nil
}

func (w *Client) url(path string) string {
	return fmt.Sprintf("http://%s:%v%s", w.host, w.port, path)
}

func mappingPath(stubMappingID string) string {
	return mappingsPath + "/" + url.PathEscape(stubMappingID)
}
