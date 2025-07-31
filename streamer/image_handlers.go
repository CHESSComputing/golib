package streamer

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// MakeOneImageReaderHandler provides one image reader handler, i.e. it reads single image from images area
func MakeOneImageReaderHandler(dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reader, err := NewImageReader("images")
		if err != nil {
			c.String(500, "init error")
			return
		}
		indexStr := c.Param("index")
		index, err := strconv.Atoi(indexStr)
		if err != nil || index < 0 {
			c.String(400, "invalid index")
			return
		}
		chunk, err := reader.ReadChunk(index)

		if err != nil {
			c.String(404, "no such image")
			return
		}
		c.Data(200, chunk.ContentType, chunk.Data)
	}
}

// MakeImageReaderHandler provides image handler based on multipart image streamer
func MakeImageReaderHandler(dir string) gin.HandlerFunc {
	return func(c *gin.Context) {
		reader, err := NewImageReader(dir)
		if err != nil {
			c.String(500, "Failed to initialize image reader: %v", err)
			return
		}
		GinMultipartImageStreamHandler(reader)(c)
	}
}

// GinMultipartImageStreamHandler provides multipart image streamer handler for dealing with multiple images in a stream
func GinMultipartImageStreamHandler(reader *ImageReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "multipart/x-mixed-replace; boundary=FRAME")

		writer := c.Writer
		for idx := 0; idx < len(reader.files); idx++ {
			chunk, err := reader.ReadChunk(idx)
			if err != nil {
				break
			}
			fmt.Fprintf(writer, "--FRAME\r\n")
			fmt.Fprintf(writer, "Content-Type: %s\r\n", chunk.ContentType)
			fmt.Fprintf(writer, "Content-Length: %d\r\n\r\n", len(chunk.Data))
			writer.Write(chunk.Data)
			fmt.Fprintf(writer, "\r\n")
			writer.Flush()
			time.Sleep(200 * time.Millisecond)
		}
	}
}
