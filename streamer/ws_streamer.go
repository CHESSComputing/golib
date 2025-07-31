package streamer

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// WebSocketStreamHandler streams binary data chunks over WebSocket
func WebSocketStreamHandler(reader BinaryReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		chunkSize := 1
		if val := c.Query("chunk"); val != "" {
			if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
				chunkSize = parsed
			}
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to upgrade to WebSocket: %v", err)
			return
		}
		defer conn.Close()

		for {
			chunk, err := reader.ReadChunk(chunkSize)
			if err != nil {
				if err == io.EOF {
					conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "done"))
					break
				}
				conn.WriteMessage(websocket.TextMessage, []byte("error: "+err.Error()))
				break
			}
			err = conn.WriteMessage(websocket.BinaryMessage, chunk.Data)
			if err != nil {
				break
			}
		}
	}
}

