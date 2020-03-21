package wiremock_test

import (
	"testing"

	"github.com/javiyt/wiremock-golang-client/pkg/wiremock"
)

func TestWiremockClient_Mappings(t *testing.T) {
	wc := wiremock.NewWiremockClient("localhost", 8080)
	wc.Mappings()
}
