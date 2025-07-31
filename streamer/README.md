# streamer

`streamer` is a Go library for streaming binary data (such as images,
scientific data, etc.) via HTTP endpoints. It supports chunked binary
streaming, WebSocket streaming, and zip packaging. The library is designed to
be extensible for machine learning, data visualization, or large dataset
delivery scenarios.

---

## Features

- Stream binary data over HTTP in customizable chunk sizes
- WebSocket support for real-time binary streaming
- ZIP archive generation on-the-fly
- Pluggable readers for different data sources (e.g., image directories, `.npy` files)
- Gin-based HTTP handlers
- Simple interface (`BinaryReader`) for implementing custom data readers
- Includes unit tests for all core components

---

## Installation

```
go get github.com/CHESSComputing/golib/streamer
````

Make sure to import it in your code:

```
import "github.com/CHESSComputing/golib/streamer"
```

---

## Usage

### Define a Reader

To use the library, implement the `BinaryReader` interface:

```
type BinaryReader interface {
    ReadChunk(chunkSize int) (*Chunk, error)
    Reset() error
}
```

A `Chunk` includes:

```
type Chunk struct {
    ContentType string
    Data        []byte
}
```

### Built-in Readers

* `ImageReader` reads images from a directory.
* `NPYReader` reads NumPy `.npy` files.

Example:

```
reader, err := streamer.NewImageReader("/path/to/images")
```

---

## HTTP Streaming Handlers

### Chunked Binary Stream

```
router.GET("/stream", streamer.GinBinaryStreamHandler(reader))
```

Query parameter `chunk` can control the chunk size (default: 1).

### ZIP Archive Download

```
router.GET("/zip", func(c *gin.Context) {
    streamer.StreamAsZip(c, reader, "bundle.zip")
})
```

Each chunk will be a separate file in the resulting ZIP archive.

### WebSocket Streaming

```
router.GET("/ws", func(c *gin.Context) {
    streamer.WebSocketStreamer(c.Writer, c.Request, reader)
})
```

Clients will receive binary messages per chunk.

---

## License

MIT License. See [LICENSE](../LICENSE) for details.

---

### Examples
Here is an example how to use streamer library with Gin web framework:
```
package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/chesscomputing/golib/streamer"
)

func main() {
	r := gin.Default()

	imgReader, err := streamer.NewImageReader("images")
	if err != nil {
		log.Fatal(err)
	}
	npyReader, err := streamer.NewNPYReader("npy")
	if err != nil {
		log.Fatal(err)
	}

	//	r.GET("/stream/images", streamer.GinBinaryStreamHandler(imgReader))
	r.GET("/stream/images", streamer.MakeImageReaderHandler("images"))
	r.GET("/image/:index", streamer.MakeOneImageReaderHandler("images"))
	r.GET("/stream/numpy", streamer.GinBinaryStreamHandler(npyReader))
	r.GET("/stream/images.zip", func(c *gin.Context) {
		reader, _ := streamer.NewImageReader("images")
		streamer.StreamAsZip(c, reader, "images_bundle.zip")
	})

	// New WebSocket Endpoints
	r.GET("/ws/images", streamer.WebSocketStreamHandler(imgReader))
	r.GET("/ws/numpy", streamer.WebSocketStreamHandler(npyReader))

	log.Println("Listening on http://localhost:8080")
	r.Run(":8080")
}

```
