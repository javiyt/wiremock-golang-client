[![CI](https://github.com/javiyt/wiremock-golang-client/actions/workflows/ci.yml/badge.svg)](https://github.com/javiyt/wiremock-golang-client/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/javiyt/wiremock-golang-client/branch/master/graph/badge.svg)](https://codecov.io/gh/javiyt/wiremock-golang-client)

# wiremock-golang-client

Go client for the [WireMock](https://wiremock.org/) admin API.

The package supports the WireMock stub mappings admin API, including listing,
creating, updating, deleting, importing, persisting, resetting, metadata
filtering, and unmatched stub operations.

## Installation

```bash
go get github.com/javiyt/wiremock-golang-client
```

## Usage

```go
package main

import (
 "fmt"
 "log"
 "net/http"

 "github.com/javiyt/wiremock-golang-client/pkg/wiremock"
)

func main() {
 client := wiremock.NewWireMockClient("localhost", 8080, nil)

 mapping, err := client.SaveMapping(wiremock.Mappings{
  Name: "hello mapping",
  Request: wiremock.Request{
   Method: http.MethodGet,
   URL:    "/hello",
  },
  Response: wiremock.Response{
   Status: http.StatusOK,
   Body:   "Hello World!!",
  },
 })
 if err != nil {
  log.Fatal(err)
 }

 fmt.Println("saved mapping:", mapping.ID)

 mappings, err := client.Mappings()
 if err != nil {
  log.Fatal(err)
 }

 fmt.Println("registered mappings:", mappings.Meta.Total)

 limit := uint(10)
 offset := uint(0)
 firstPage, err := client.MappingsWithOptions(wiremock.MappingsOptions{
  Limit:  &limit,
  Offset: &offset,
 })
 if err != nil {
  log.Fatal(err)
 }

 fmt.Println("first page mappings:", len(firstPage.Mappings))

 updated, err := client.UpdateMapping(mapping.ID, wiremock.Mappings{
  Name: "updated hello mapping",
  Request: wiremock.Request{
   Method: http.MethodGet,
   URL:    "/hello",
  },
  Response: wiremock.Response{
   Status: http.StatusOK,
   Body:   "Hello again!!",
  },
  Metadata: map[string]string{"tag": "examples"},
 })
 if err != nil {
  log.Fatal(err)
 }

 fmt.Println("updated mapping:", updated.ID)

 matched, err := client.FindMappingsByMetadata(wiremock.MetadataFilter{
  "matchesJsonPath": map[string]any{
   "expression": "$.tag",
   "equalTo":    "examples",
  },
 })
 if err != nil {
  log.Fatal(err)
 }

 fmt.Println("matched mappings:", matched.Meta.Total)

 if err := client.PersistMappings(); err != nil {
  log.Fatal(err)
 }

 if err := client.DeleteMapping(mapping.ID); err != nil {
  log.Fatal(err)
 }
}
```

Pass a custom `*http.Client` to `NewWireMockClient` when you need custom
timeouts, transports, or test doubles. Passing `nil` uses a default HTTP client.

## Development

Run the unit and functional test suite:

```bash
make test
```

The functional tests expect WireMock on `localhost:8080`. The `make test`
target starts the Docker Compose service defined in `docker-compose.yml`, runs
the test script, and stops the service afterwards.

Useful commands:

```bash
make format
make linter
make run-test
```

## License

This project is licensed under the Apache License 2.0. See [LICENSE](LICENSE).
