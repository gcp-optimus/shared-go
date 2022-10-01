package modules

import (
	"context"
	"sync"

	"cloud.google.com/go/datastore"
	"github.com/gcp-optimus/shared-go/logger"
	"go.uber.org/zap"
)

// global (instance-wide) scope for reusable when wakeup
var datastoreClientOnce sync.Once

type datastoreImpl struct {
	client    *datastore.Client
	projectId string
}

func GetDatastoreModule(projectId string) DatastoreRepository {
	return &datastoreImpl{projectId: projectId}
}

func (entry *datastoreImpl) GetDb() *datastore.Client {
	entry.sync()
	return entry.client
}

func (entry *datastoreImpl) sync() {
	// You may wish to add different checks to see if the client is needed for this request.
	datastoreClientOnce.Do(func() {
		// Pre-declare an err variable to avoid shadowing client.
		var err error
		entry.client, err = datastore.NewClient(context.Background(), entry.projectId)
		if err != nil {
			logger.Critical("can't get datastore client",
				zap.String("func", "datastore.NewClient"),
				zap.String("projectId", entry.projectId),
				zap.Error(err),
			)
		}
	})
}
