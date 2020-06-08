package wiremock_golang_client_test

import (
	"github.com/javiyt/wiremock-golang-client/pkg/wiremock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWireMockClient(t *testing.T) {
	wc := wiremock.NewWireMockClient("localhost", 8080, nil)
	mapping, err := wc.Mappings()

	assert.NoError(t, err)
	assert.Equal(t, uint(2), mapping.Meta.Total)
	assert.Len(t, mapping.Mappings, 2)
}
