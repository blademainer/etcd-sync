package tests

import (
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/blademainer/etcd-sync/etcdsync"
	log "github.com/blademainer/etcd-sync/logger"
	"sync"
	"testing"
)

var logger = log.Log

type (
	EtcdSync struct {
		Master  *etcdsync.MasterConfig         `yaml:"master" json:"master"`
		Slaves  []*EtcdConfig                  `yaml:"slaves" json:"slaves"`
		Options []etcdsync.EtcdOptionInterface `yaml:"-" json:"-"`
		wg      *sync.WaitGroup                `yaml:"-" json:"-"`
	}

	EtcdConfig struct {
		Endpoints      []string             `yaml:"endpoints" json:"endpoints"`
		TimeoutSeconds int64                `yaml:"TimeoutSeconds" json:"timeout_seconds"`
		client         *clientv3.Client     `yaml:"-" json:"-"`
		channel        chan *clientv3.Event `yaml:"-" json:"-"`
		closeChan      chan bool            `yaml:"-" json:"-"`
		disconnected   bool                 `yaml:"-" json:"-"`
		sync.Mutex
	}
)

func TestPtr(t *testing.T) {
	etcdSync := &EtcdSync{}
	etcdSync.Slaves = make([]*EtcdConfig, 1)
	etcdSync.Slaves[0] = &EtcdConfig{}
	etcdSync.Slaves[0].disconnected = true
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for _, e := range etcdSync.Slaves {
			e.initEtcdClient()
			fmt.Println("disconnected: ", e.disconnected)
		}
		wg.Done()
	}()
	wg.Wait()

}

func (e *EtcdConfig) initEtcdClient() {
	e.disconnected = false
}
