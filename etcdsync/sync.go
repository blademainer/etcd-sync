package etcdsync

import (
	"github.com/coreos/etcd/clientv3"
	log "github.com/blademainer/etcd-sync/logger"
	sync2 "sync"
	"time"
)

var logger = log.Log

func Init(sync *EtcdSync) {
	config := sync.Master
	config.initEtcdClient()

	for _, e := range sync.Slaves {
		e.initEtcdClient()
	}

	options := make([]EtcdOptionInterface, 0)

	for _, e := range sync.Master.FixedConfigs {
		options = append(options, e)
	}
	for _, e := range sync.Master.PrefixConfigs {
		options = append(options, e)
	}
	for _, e := range sync.Master.RangeConfigs {
		options = append(options, e)
	}
	sync.Options = options

	wg := &sync2.WaitGroup{}
	sync.wg = wg

	config.closeChan = make(chan bool, len(options))
}

func (e *EtcdConfig) initEtcdClient() {
	e.Lock()
	defer e.Unlock()
	logger.Infof("Connecting etcd: %v", e.Endpoints)
	c := clientv3.Config{Endpoints: e.Endpoints, DialTimeout: time.Duration(e.TimeoutSeconds) * time.Second}

	cli, err := clientv3.New(c)
	if err != nil {
		logger.Errorf("Error when connection: %s error: %s \n", e.Endpoints, err.Error())
		e.disconnected = true
	} else if e.client != nil {
		e.client.Close()
	} else {
		e.disconnected = false
	}
	e.client = cli
	if e.channel == nil {
		e.channel = make(chan *clientv3.Event, 1024)
	} else if e.closeChan == nil {
		e.closeChan = make(chan bool, 1)
	}
	logger.Infof("Inited etcdClient: %v", cli)
}

func (sync *EtcdSync) Start() {
	Init(sync)
	logger.Infof("Start... SyncOnce...")
	sync.SyncOnce()
	logger.Infof("Done SyncOnce.")

	logger.Infof("Starting watcher...")
	sync.StartWatcher()
	logger.Infof("Exited watcher...")
}
