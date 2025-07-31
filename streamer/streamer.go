package streamer

import (
	"io"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Chunk defines chunk data-structure
type Chunk struct {
	ContentType string
	Data        []byte
}

// BinaryReader provides binary reader
type BinaryReader interface {
	ReadChunk(chunkSize int) (*Chunk, error)
	Reset() error
}

// GinBinaryStreamHandler provides binary streamer handler function for gin web framework
func GinBinaryStreamHandler(reader BinaryReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		chunkSize := 1
		if val := c.Query("chunk"); val != "" {
			if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
				chunkSize = parsed
			}
		}

		firstChunk, err := reader.ReadChunk(chunkSize)
		if err != nil {
			if err == io.EOF {
				c.Status(204)
				return
			}
			c.String(500, "Error: %v", err)
			return
		}

		// Set headers BEFORE writing
		c.Header("Content-Type", firstChunk.ContentType)
		c.Header("Transfer-Encoding", "chunked")
		c.Status(200)

		// Write first chunk
		c.Writer.Write(firstChunk.Data)
		c.Writer.Flush()

		// Write remaining chunks
		for {
			chunk, err := reader.ReadChunk(chunkSize)
			if err != nil {
				if err == io.EOF {
					break
				}
				c.String(500, "Error: %v", err)
				return
			}
			c.Writer.Write(chunk.Data)
			c.Writer.Flush()
		}
	}
}
