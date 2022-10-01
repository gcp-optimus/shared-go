package modules

import (
	"net/http"
)

// global (instance-wide) scope for reusable when wakeup
type httpImpl struct {
	client *http.Client
}

func GetHttpModule() HttpClient {
	return &httpImpl{client: http.DefaultClient}
}

func (entry *httpImpl) GetClient() *http.Client {
	return entry.client
}
