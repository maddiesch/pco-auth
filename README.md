# Planning Center Authentication

Planning Center Authentication can be used to authenticate an Oauth application with the [Planning Center API](https://developer.planning.center/docs/).

The `Perform` method is the main method for authenticating. It will open a browser, and start an HTTP server to listen for the callback URL on the port specified.

```go
package main

import (
  "fmt"
  "log"

  auth "github.com/maddiesch/pco-auth"
)

func main() {
  input := &auth.PerformInput{
    Port:         8080,
    ClientID:     "0000000000000000000000000000000000000000000000000000000000000000",
    ClientSecret: "0000000000000000000000000000000000000000000000000000000000000000",
    Scopes:       []string{"people"}
  }

  output, err := auth.Perform(&input)
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(output)
}
```

## Example

[See the example CLI app here](/example)

## Note

Any example & test Client ID & Client Secret values are fake. You will need to create your own at [api.planningcenteronline.com/oauth/applications](https://api.planningcenteronline.com/oauth/applications).
