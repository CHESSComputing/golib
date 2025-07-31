package streamer

import (
	"archive/zip"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

// StreamAsZip provides zip streamer
func StreamAsZip(c *gin.Context, reader BinaryReader, zipName string) {
	c.Writer.Header().Set("Content-Type", "application/zip")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+zipName)

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	chunkNum := 0
	for {
		chunk, err := reader.ReadChunk(chunkNum)
		if err == io.EOF {
			break
		}
		if err != nil {
			c.String(500, "read error: %v", err)
			return
		}

		fw, err := zipWriter.Create(fmt.Sprintf("chunk_%03d", chunkNum))
		if err != nil {
			c.String(500, "zip error: %v", err)
			return
		}

		_, err = fw.Write(chunk.Data)
		if err != nil {
			c.String(500, "write error: %v", err)
			return
		}

		chunkNum++
	}

}
