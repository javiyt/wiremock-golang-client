package wiremock_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/javiyt/wiremock-golang-client/pkg/wiremock"

	_ "github.com/mailru/easyjson/gen"
)

func TestWiremockClient_Mappings(t *testing.T) {
	wc := wiremock.NewWiremockClient("localhost", 8080)
	mapping, err := wc.Mappings()

	assert.NoError(t, err)
	assert.Equal(t, uint(2), mapping.Meta.Total)
	assert.Len(t, mapping.Mappings, 2)
}
