package modules

import (
	"net/http"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/datastore"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/websocket"
)

type Boot interface {
	WaitForTerminateSignal()
}

type DatastoreRepository interface {
	GetDb() *datastore.Client
}

type FiberServer interface {
	GetServer() *fiber.App
	Serve()
}

type CloudTasksClient interface {
	GetClient() *cloudtasks.Client
}

type HttpClient interface {
	GetClient() *http.Client
}

type WebSocketClient interface {
	Dial() (*websocket.Conn, error)
}

func GetEnv(key, default_ string) string {
	env := os.Getenv(key)
	if env == "" {
		return default_
	}

	return env
}
