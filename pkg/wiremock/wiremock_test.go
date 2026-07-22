package wiremock_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/javiyt/wiremock-golang-client/pkg/wiremock"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (r roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

type errorReadCloser struct{}

func (e errorReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("error reading mappings")
}

func (e errorReadCloser) Close() error {
	return nil
}

func TestClient_Mappings(t *testing.T) {
	t.Run("it should fail when not possible to get mappings", func(t *testing.T) {
		wClient := wiremock.NewWireMockClient("localhost", 8000, &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("error accessing mappings")
			}),
		})
		_, err := wClient.Mappings()

		require.ErrorContains(t, err, "error accessing mappings")
	})

	t.Run("it should fail when API response is not ok", func(t *testing.T) {
		wClient := wiremock.NewWireMockClient("localhost", 8000, &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			}),
		})

		_, err := wClient.Mappings()

		require.EqualError(t, err, "error got from API, status code: 500")
	})

	t.Run("it should fail when response body cannot be read", func(t *testing.T) {
		wClient := wiremock.NewWireMockClient("localhost", 8000, &http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       errorReadCloser{},
				}, nil
			}),
		})

		_, err := wClient.Mappings()

		require.EqualError(t, err, "error reading response body: error reading mappings")
	})

	t.Run("it should fail when response body is empty", func(t *testing.T) {
		wClient := wiremockClientFromServer(t, "")

		_, err := wClient.Mappings()

		require.ErrorContains(t, err, "error unmarshaling response")
	})

	t.Run("it should fail when response body is a malformed JSON", func(t *testing.T) {
		wClient := wiremockClientFromServer(t, "{\"mappings")

		_, err := wClient.Mappings()

		require.ErrorContains(t, err, "error unmarshaling response")
	})

	t.Run("it should return a registered mapping", func(t *testing.T) {
		wClient := wiremockClientFromServer(t, `{
			"mappings": [{
				"id": "012e3261-3398-46da-9811-deb02de35872",
				"uuid": "012e3261-3398-46da-9811-deb02de35872",
				"name": "hello mapping",
				"request": {
					"url": "/hello",
					"urlPathPattern": "/hello/.*",
					"method": "GET",
					"queryParameters": {"name": "javi"},
					"headers": {"Accept": "text/plain"},
					"cookies": {"session": "abc"}
				},
				"response": {
					"status": 200,
					"statusMessage": "OK",
					"body": "Hello World!!",
					"headers": {"Content-Type": "text/plain"},
					"fixedDelayMilliseconds": 10,
					"transformers": ["response-template"]
				},
				"persistent": true,
				"priority": 1,
				"scenarioName": "hello scenario",
				"requiredScenarioState": "Started",
				"newScenarioState": "Completed",
				"metadata": {"owner": "tests"}
			}],
			"meta": {"total": 1}
		}`)

		mappings, err := wClient.Mappings()

		require.NoError(t, err)
		require.Equal(t, uint(1), mappings.Meta.Total)
		require.Len(t, mappings.Mappings, 1)
		require.Equal(t, "012e3261-3398-46da-9811-deb02de35872", mappings.Mappings[0].ID)
		require.Equal(t, "012e3261-3398-46da-9811-deb02de35872", mappings.Mappings[0].UUID)
		require.Equal(t, "hello mapping", mappings.Mappings[0].Name)
		require.Equal(t, "/hello", mappings.Mappings[0].Request.URL)
		require.Equal(t, "/hello/.*", mappings.Mappings[0].Request.URLPathPattern)
		require.Equal(t, "GET", mappings.Mappings[0].Request.Method)
		require.Equal(t, "javi", mappings.Mappings[0].Request.QueryParameters["name"])
		require.Equal(t, "text/plain", mappings.Mappings[0].Request.Headers["Accept"])
		require.Equal(t, "abc", mappings.Mappings[0].Request.Cookies["session"])
		require.Equal(t, uint(200), mappings.Mappings[0].Response.Status)
		require.Equal(t, "OK", mappings.Mappings[0].Response.StatusMessage)
		require.Equal(t, "Hello World!!", mappings.Mappings[0].Response.Body)
		require.Equal(t, "text/plain", mappings.Mappings[0].Response.Headers["Content-Type"])
		require.Equal(t, uint(10), mappings.Mappings[0].Response.FixedDelayMilliseconds)
		require.Equal(t, []string{"response-template"}, mappings.Mappings[0].Response.Transformers)
		require.True(t, mappings.Mappings[0].Persistent)
		require.Equal(t, uint(1), mappings.Mappings[0].Priority)
		require.Equal(t, "hello scenario", mappings.Mappings[0].ScenarioName)
		require.Equal(t, "Started", mappings.Mappings[0].RequiredScenarioState)
		require.Equal(t, "Completed", mappings.Mappings[0].NewScenarioState)
		require.Equal(t, "tests", mappings.Mappings[0].Metadata["owner"])
	})
}

func wiremockClientFromServer(t *testing.T, body string) *wiremock.Client {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodGet, req.Method)
		require.Equal(t, "/__admin/mappings", req.URL.Path)

		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(body))
		require.NoError(t, err)
	}))
	t.Cleanup(server.Close)

	host, portValue, found := strings.Cut(strings.TrimPrefix(server.URL, "http://"), ":")
	require.True(t, found)

	port, err := strconv.ParseUint(portValue, 10, 0)
	require.NoError(t, err)

	return wiremock.NewWireMockClient(host, uint(port), nil)
}
