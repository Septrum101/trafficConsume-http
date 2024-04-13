package app

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/Septrum101/trafficConsume-http/common"
	"github.com/Septrum101/trafficConsume-http/infra"
)

func New() *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &Client{
		Client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: tr,
		},
		bufPool: &sync.Pool{
			New: func() interface{} {
				b := make([]byte, 8192)
				return &b
			},
		},
		stopChan: make(chan bool),
	}
}

func (c *Client) download(downloadUrl string) {
	resp, err := c.Get(downloadUrl)
	if err != nil {
		log.Debugln(err)
		return
	}
	defer resp.Body.Close()

	buf := c.bufPool.Get().(*[]byte)
	defer func() {
		c.bufPool.Put(buf)
	}()

	for {
		n, _ := resp.Body.Read(*buf)
		if n == 0 {
			return
		}

		c.readBytes.Add(int64(n))
	}
}

func (c *Client) monitor() {
	oldBytes := int64(0)
	old := time.Now()
	for n := range time.Tick(time.Second * 30) {
		newBytes := c.readBytes.Load()
		speed := (newBytes - oldBytes) * 1000 / n.Sub(old).Milliseconds()
		oldBytes = newBytes
		old = n
		log.Infof("Throughput: %s, speed: %s/s", infra.ByteCountIEC(newBytes), infra.ByteCountIEC(speed))
	}
}

func (c *Client) Start() error {
	log.Info("Starting..")
	urls, err := common.GetDownloadUrl()
	if err != nil {
		return err
	}

	go c.monitor()

	go func() {
		c.wg.Add(1)
		defer c.wg.Done()

		worker := make(chan bool, 32)
		for {
			select {
			case <-c.stopChan:
				log.Info("Closing...")
				return
			default:
				worker <- true
				go func() {
					defer func() {
						<-worker
					}()
					r := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(urls))
					c.download(urls[r])
				}()
			}
		}
	}()

	return nil
}

func (c *Client) Close() {
	close(c.stopChan)
	c.wg.Wait()
}
