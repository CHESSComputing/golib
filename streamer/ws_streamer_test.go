package streamer

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// TestWebSocketStreamHandler provides unit test for web socket stream handler
func TestWebSocketStreamHandler(t *testing.T) {
	dir := t.TempDir()
	setupTestImages(t, dir)

	reader, err := NewImageReader(dir)
	if err != nil {
		t.Fatal(err)
	}

	r := gin.Default()
	r.GET("/ws", WebSocketStreamHandler(reader))

	srv := httptest.NewServer(r)
	defer srv.Close()

	wsURL := "ws" + srv.URL[4:] + "/ws"
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect websocket: %v", err)
	}
	defer conn.Close()

	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("error reading message: %v", err)
	}

	if len(data) == 0 {
		t.Error("received empty websocket message")
	}
}
