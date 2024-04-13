package app

import (
	"net/http"
	"sync"
	"sync/atomic"
)

type Client struct {
	*http.Client

	readBytes atomic.Int64
	bufPool   *sync.Pool
	wg        sync.WaitGroup
	stopChan  chan bool
}
