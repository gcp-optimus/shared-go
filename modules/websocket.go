package modules

import (
	"fmt"

	"github.com/gcp-optimus/shared-go/logger"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

// global (instance-wide) scope for reusable when wakeup
type webSocketImpl struct {
	url string
}

func GetWebSocketModule(schema, host, port string) WebSocketClient {
	return &webSocketImpl{url: fmt.Sprintf("%s://%s:%s", schema, host, port)}
}

func (entry *webSocketImpl) Dial() (*websocket.Conn, error) {
	conn, err := websocket.Dial(entry.url, "", "")
	if err != nil {
		logger.Error("can't get dial web socket",
			zap.String("func", "websocket.Dial"),
			zap.String("url", entry.url),
			zap.Error(err),
		)
		return nil, err
	}

	return conn, nil
}
