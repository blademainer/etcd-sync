package etcdsync

import (
	"context"
	"github.com/coreos/etcd/clientv3"
)

func (sync *EtcdSync) StartWatcher() {
	// Start master watcher
	sync.startMaster()
	sync.startSlaves()
	sync.wg.Wait()
}

func (sync *EtcdSync) CloseWatcher() {
	sync.closeMaster()
	sync.closeSlave()
}

func (sync *EtcdSync) startMaster() {
	logger.Infof("Options: %v size: %d", sync.Options, len(sync.Options))
	sync.wg.Add(len(sync.Options))
	for _, e := range sync.Options {
		go func(opt EtcdOptionInterface) {
			defer logger.Warnf("Exit!!!")
			defer sync.wg.Done()
			watch := sync.Master.client.Watch(context.Background(), opt.GetKey(), opt.OpOptions())
			for {
				select {
				case ev := <-watch:
					if logger.IsInfoEnabled() {
						logger.Debugf("Receive response: %v", ev)
					}
					sync.processWatch(ev)
				case <-sync.Master.closeChan:
					logger.Warnf("Closing master...")
					return
				}
			}
		}(e)
	}
}

func (sync *EtcdSync) closeMaster() {
	for range sync.Options {
		sync.Master.closeChan <- true
		sync.Master.closed = true
	}
}

func (sync *EtcdSync) processWatch(watchResponse clientv3.WatchResponse) {
	// send data to all slave channel
	for _, e := range watchResponse.Events {
		for _, slave := range sync.Slaves {
			if !slave.disconnected {
				slave.channel <- e
				logger.Debugf("Succeed put event: %v to slave channel: %v", e, *slave)
			} else {
				logger.Warnf("Slave is disconnected!! so ignore push to channel. slave: %v", *slave)
			}
		}
	}
}
