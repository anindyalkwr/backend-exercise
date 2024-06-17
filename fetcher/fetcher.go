package fetcher

import (
	"net/http"
	"sync"
)

const appIDHeader = "app-id"

type Fetcher interface {
	Fetch(client *http.Client, appId string, page int, wg *sync.WaitGroup)
}
