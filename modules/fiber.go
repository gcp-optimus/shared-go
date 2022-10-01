package modules

import (
	"fmt"

	"github.com/gcp-optimus/shared-go/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type fiberImpl struct {
	server *fiber.App
	port   string
}

func GetFiberModule(name string, port string) FiberServer {
	app := fiber.New(fiber.Config{
		ReduceMemoryUsage:     true,
		Prefork:               false,
		CaseSensitive:         true,
		StrictRouting:         true,
		AppName:               name,
		DisableStartupMessage: true,
	})

	app.Use(zapLogger)

	return &fiberImpl{server: app, port: port}
}

func (entry *fiberImpl) GetServer() *fiber.App {
	return entry.server
}

func (entry *fiberImpl) Serve() {
	port := fmt.Sprintf(":%s", entry.port)

	go func(port string) {
		if err := entry.server.Listen(port); err != nil {
			logger.Critical("can't listen by fiber server",
				zap.String("func", "fiber.Listen"),
				zap.String("port", port),
				zap.Error(err),
			)
		}
	}(port)
}

func zapLogger(c *fiber.Ctx) error {
	err := c.Next()
	req := c.Request()
	res := c.Response()
	code := res.StatusCode()

	if err != nil {
		logger.Error("unexpected error",
			zap.String("Method", c.Method()),
			zap.ByteString("URI", req.URI().FullURI()),
			zap.ByteString("Header", req.Header.Header()),
			zap.ByteString("Body", req.Body()),
			zap.Error(err))
	} else if code >= 500 {
		// internal server error
		logger.Error("internal server error",
			zap.String("Method", c.Method()),
			zap.ByteString("URI", req.URI().FullURI()),
			zap.ByteString("Header", req.Header.Header()),
			zap.ByteString("Body", req.Body()),
			zap.Int("Code", code))
	} else if code >= 400 {
		// client error
		logger.Warn("client error",
			zap.String("Method", c.Method()),
			zap.ByteString("URI", req.URI().FullURI()),
			zap.ByteString("Header", req.Header.Header()),
			zap.ByteString("Body", req.Body()),
			zap.Int("Code", code))
	} else {
		logger.Info("",
			zap.String("Method", c.Method()),
			zap.ByteString("URI", req.URI().FullURI()),
			zap.Int("Code", code))
	}

	return err
}
