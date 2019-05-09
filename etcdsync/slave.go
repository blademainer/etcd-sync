package etcdsync

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"time"
)

func (sync *EtcdSync) startSlaves() {
	sync.wg.Add(len(sync.Slaves))
	for _, slave := range sync.Slaves {
		go func(s *EtcdConfig) {
			logger.Infof("Started slave processor... %v", s)
			defer sync.wg.Done()
			sync.startSlaveProcessor(s)
		}(slave)
	}
	sync.startCheckSlaveStatus()
}

func (sync *EtcdSync) startSlaveProcessor(slave *EtcdConfig) {
	for {
		select {
		case e := <-slave.channel:
			slave.processSlaveEvent(e)
		case <-slave.closeChan:
			return
		}
	}
}

func (slave *EtcdConfig) processSlaveEvent(e *clientv3.Event) {
	if slave.disconnected {
		logger.Errorf("Slave: %v is disconnected! so break process channel!", *slave)
		return
	}
	key := string(e.Kv.Key)
	value := string(e.Kv.Value)
	timeoutCtx, _ := context.WithTimeout(context.TODO(), time.Duration(slave.TimeoutSeconds)*time.Second)
	switch e.Type {
	case mvccpb.PUT:
		logger.Debugf("Putting key: %s value: %s to slave: %v", key, value, *slave)
		if response, err := slave.client.Put(timeoutCtx, key, value); err != nil {
			logger.Errorf("Failed to putting data on slave: %v error: %v, key: %s value: %s", *slave, err, key, value)
			slave.isDisconnected(err)
		} else {
			logger.Infof("Succeed put data on slave: %v with response: %v, key: %s value: %s", *slave, response, key, value)
		}
	case mvccpb.DELETE:
		logger.Warnf("Deleting key: [%s] on client: %v", key, *slave)
		if response, err := slave.client.Delete(timeoutCtx, key); err != nil {
			logger.Errorf("Failed to delete data on slave: %v error: %v, key: %s value: %s", *slave, err, key, value)
			slave.isDisconnected(err)
		} else {
			logger.Infof("Deleted data on slave: %v with response: %v, key: %s value: %s", *slave, response, key, value)
		}
	}
}

func (slave *EtcdConfig) isDisconnected(err error) {
	if err != context.DeadlineExceeded {
		return
	}
	slave.disconnected = true
}

func (sync *EtcdSync) closeSlave() {
	for _, slave := range sync.Slaves {
		slave.closeChan <- true
		slave.disconnected = true
	}
}

// -------- check slave status
func (sync *EtcdSync) startCheckSlaveStatus() {
	tick := time.Tick(10 * time.Second)

	go func() {
		for range tick {
			sync.checkSlaveStatus()
		}
	}()
}

func (sync *EtcdSync) checkSlaveStatus() {
	for _, slave := range sync.Slaves {
		if slave.disconnected {
			for _, endpoint := range slave.Endpoints {
				timeoutCtx, _ := context.WithTimeout(context.TODO(), time.Duration(slave.TimeoutSeconds)*time.Second)
				response, err := slave.client.Status(timeoutCtx, endpoint)
				if err != nil{
					logger.Errorf("Failed to check status, error: %s, endpoint: %s", err.Error(), endpoint)
				} else if response != nil {
					logger.Infof("Got status of endpoint: %s response: %v. And reset status to connected.", endpoint, response)
					slave.disconnected = false
				}
			}
			logger.Warnf("Slave is closed: %v try to restart.", *slave)
			slave.initEtcdClient()
			if !slave.disconnected {
				logger.Warnf("Reconnected...")
				sync.SyncMasterToSlave(slave)
			} else {
				logger.Errorf("Failed to reconnect etcd: %v", slave.Endpoints)
			}
		}
	}
}
