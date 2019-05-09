package etcdsync

import (
	"context"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

func (sync *EtcdSync) SyncOnce() {
	sync.syncMasterToClient(sync.Options, sync.Slaves...)
}

func (sync *EtcdSync) SyncMasterToSlave(slaves ...*EtcdConfig) {
	sync.syncMasterToClient(sync.Options, slaves...)
}

func (sync *EtcdSync) syncMasterToClient(optionInterfaces []EtcdOptionInterface, slaves ...*EtcdConfig) {
	for _, e := range optionInterfaces {
		if response, err := sync.Master.client.Get(context.TODO(), e.GetKey(), e.OpOptions()); err != nil {
			logger.Errorf("Failed to get key: %s error: %s", e.GetKey(), err.Error())
		} else {
			values := response.Kvs
			sync.putKvsToSlaves(values, slaves...)
		}
	}
}

func (sync *EtcdSync) putKvsToSlaves(kvs []*mvccpb.KeyValue, slaves ...*EtcdConfig) {
	for _, slave := range slaves {
		if slave.disconnected {
			logger.Warnf("Client: %v is closed!!! so ignored to update config.", slave.Endpoints)
			continue
		}
		for _, kv := range kvs {
			if response, err := slave.client.Put(context.TODO(), string(kv.Key), string(kv.Value)); err != nil {
				logger.Errorf("Failed to put kv: %v to client: %v with error: %s", kv, slave.Endpoints, err.Error())
			} else {
				logger.Infof("Succeed to put kv: %v to client: %v. Response: %v", kv, slave.Endpoints, response)
			}
		}
	}
}
