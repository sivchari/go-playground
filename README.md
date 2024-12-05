# goplayground

goplayground is a client to call Go Playground API.

## Installation

```bash
go get -u github.com/otiai10/goplayground
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/sivchari/playground"
)

func main() {
    c := playground.NewClient()
    code := `package main
import "fmt"

func main() {
    fmt.Println("Hello, playground")
}`
    res, err := c.Run([]byte(code))
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(res)
}
```

## Configuration

You can configure the url of the Go Playground API.
Following environment variables are available.

- `PLAYGROUND_FRONTEND_URL` (default: `https://go.dev/play`)
- `PLAYGROUND_BACKEND_URL` (default: `https://play.golang.org`)
    



