package modules

import (
	"context"
	"sync"

	"github.com/gcp-optimus/shared-go/logger"
	"go.uber.org/zap"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
)

// global (instance-wide) scope for reusable when wakeup
var cloudtasksClientOnce sync.Once

type cloudTasksImpl struct {
	client *cloudtasks.Client
}

func GetCloudTasksModule() CloudTasksClient {
	return &cloudTasksImpl{}
}

func (entry *cloudTasksImpl) GetClient() *cloudtasks.Client {
	entry.sync()
	return entry.client
}

func (entry *cloudTasksImpl) sync() {
	// You may wish to add different checks to see if the client is needed for this request.
	cloudtasksClientOnce.Do(func() {
		// Pre-declare an err variable to avoid shadowing client.
		var err error
		entry.client, err = cloudtasks.NewClient(context.Background())
		if err != nil {
			logger.Critical("can't get cloudtasks client",
				zap.String("func", "cloudtasks.NewClient"),
				zap.Error(err),
			)
		}
	})
}
