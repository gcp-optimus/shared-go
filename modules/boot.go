package modules

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gcp-optimus/shared-go/logger"
)

type bootImpl struct {
	startTime       time.Time
	terminateSignal chan os.Signal
}

var (
	appCtx = &bootImpl{
		startTime:       time.Now(),
		terminateSignal: make(chan os.Signal, 1),
	}
)

// Init global app context with bellow fields.
func init() {
	signal.Notify(
		appCtx.terminateSignal,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
}

func GetBoot(env string) Boot {
	logger.GlobalBuild(env == "production")
	return appCtx
}

// WaitForTerminateSignal waits for shutdown signal.
func (ctx *bootImpl) WaitForTerminateSignal() {
	logger.Info("Running Service")
	<-ctx.terminateSignal
	logger.Info("Terminated Service")
}
