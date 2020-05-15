package wiremock_test

import (
	"errors"
	"github.com/jarcoal/httpmock"
	"github.com/javiyt/wiremock-golang-client/pkg/wiremock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestClient_Mappings(t *testing.T) {
	client := &http.Client{Transport: &http.Transport{}}
	httpmock.ActivateNonDefault(client)
	defer httpmock.DeactivateAndReset()

	wClient := wiremock.NewWireMockClient("localhost", 8000, client)

	t.Run("it should fail when not possible to get mappings", func(t *testing.T) {
		httpmock.RegisterResponder(
			"GET",
			"http://localhost:8000/__admin/mappings",
			httpmock.NewErrorResponder(errors.New("error accessing mappings")),
		)

		_, err := wClient.Mappings()

		require.EqualError(t, err, "Get \"http://localhost:8000/__admin/mappings\": error accessing mappings")
	})

	t.Run("it should fail when response body is empty", func(t *testing.T) {
		httpmock.RegisterResponder(
			"GET",
			"http://localhost:8000/__admin/mappings",
			httpmock.NewStringResponder(200, ""),
		)

		_, err := wClient.Mappings()

		require.EqualError(t, err, "error unmarshaling response readObjectStart: expect { or n, but found \x00, error found in #0 byte of ...||..., bigger context ...||...")
	})

	t.Run("it should fail when response body is a malformed JSON", func(t *testing.T) {
		httpmock.RegisterResponder(
			"GET",
			"http://localhost:8000/__admin/mappings",
			httpmock.NewStringResponder(200, "{\"mappings"),
		)

		_, err := wClient.Mappings()

		require.EqualError(t, err, "error unmarshaling response wiremock.Mapping.readFieldHash: incomplete field name, error found in #10 byte of ...|{\"mappings|..., bigger context ...|{\"mappings|...")
	})

	t.Run("it should return a registered mapping", func(t *testing.T) {
		responder := httpmock.NewStringResponder(200, "{\"mappings\":[{\"id\":\"012e3261-3398-46da-9811-deb02de35872\",\"request\":{\"url\":\"/hello\",\"method\":\"GET\"},\"response\":{\"status\":200,\"body\":\"Hello World!!\",\"headers\":{\"Content-Type\":\"text/plain\"}},\"uuid\":\"012e3261-3398-46da-9811-deb02de35872\"}],\"meta\":{\"total\":1}}")

		httpmock.RegisterResponder(
			"GET",
			"http://localhost:8000/__admin/mappings",
			responder,
		)

		mappings, err := wClient.Mappings()

		require.NoError(t, err)
		require.Len(t, mappings.Mappings, 1)
		require.Equal(t, "012e3261-3398-46da-9811-deb02de35872", mappings.Mappings[0].ID)
		require.Equal(t, "/hello", mappings.Mappings[0].Request.URL)
		require.Equal(t, "GET", mappings.Mappings[0].Request.Method)
		require.Equal(t, uint(200), mappings.Mappings[0].Response.Status)
		require.Equal(t, "Hello World!!", mappings.Mappings[0].Response.Body)
		require.Equal(t, "text/plain", mappings.Mappings[0].Response.Headers["Content-Type"])
	})
}
