# M3U

A basic golang [M3U playlist](https://en.wikipedia.org/wiki/M3U) parser library.

## Usage
```
go get github.com/jamesnetherton/m3u
```

```go
package main

import (
	"fmt"
	"github.com/jamesnetherton/m3u"
)

func main() {
	playlist, err := m3u.Parse("playlist.m3u")

	if err == nil {
		for _, track := range playlist.Tracks {
			fmt.Println("Track name = ", track.Name)
			fmt.Println("Track length = ", track.Length)
			fmt.Println("Track URI = ", track.URI)
		}
	}
}
```
