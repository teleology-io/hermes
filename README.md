# hermes
A very small wrapper around golang http requests

# Installation

```
go get github.com/teleology-io/hermes
```

# Usage 

Creating a client:

```go
package main

import (
	"fmt"

	"github.com/teleology-io/hermes"
)

func main() {
	client := hermes.Create(hermes.ClientConfiguration{
		BaseURL: "some_url_here",
		Headers: hermes.Headers{},
		Params: hermes.Params{},
		Timeout: 5, // seconds
	})
}
```

Making requests:


```go
res, err := client.Send(hermes.Request{
		Method:  hermes.POST,
		Url:     "/v1/configuration",
		Headers: hermes.Headers{}, // per-request
		Params:  hermes.Params{},  // per-request
		Data:    "some_interface",
	})
if err != nil {
  // handle err
}

// res wraps http.Response but automatically reads Body 
// and stores it on res.Data
```