package wiremock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/javiyt/wiremock-golang-client/pkg/wiremock"
)

func TestWireMockClient_Mappings(t *testing.T) {
	wc := wiremock.NewWireMockClient("localhost", 8080)
	mapping, err := wc.Mappings()

	assert.NoError(t, err)
	assert.Equal(t, uint(2), mapping.Meta.Total)
	assert.Len(t, mapping.Mappings, 2)
}
