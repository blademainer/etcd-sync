package etcdsync

import (
	"github.com/coreos/etcd/clientv3"
	"sync"
)

type (
	EtcdSync struct {
		Master  *MasterConfig         `yaml:"master" json:"master"`
		Slaves  []*EtcdConfig         `yaml:"slaves" json:"slaves"`
		Options []EtcdOptionInterface `yaml:"-" json:"-"`
		wg      *sync.WaitGroup       `yaml:"-" json:"-"`
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

	MasterConfig struct {
		FixedConfigs  []*FixedConfig  `yaml:"fixedConfigs" json:"fixed_configs"`
		PrefixConfigs []*PrefixConfig `yaml:"prefixConfigs" json:"prefix_configs"`
		RangeConfigs  []*RangeConfig  `yaml:"rangeConfigs" json:"range_configs"`
		EtcdConfig
		closeChan chan bool `yaml:"-"  json:"-"`
		closed    bool      `yaml:"-" json:"-"`
	}

	EtcdOptionInterface interface {
		OpOptions() clientv3.OpOption
		GetKey() string
	}

	FixedConfig struct {
		Key string `yaml:"key"`
	}

	PrefixConfig struct {
		Key string `yaml:"key"`
	}

	RangeConfig struct {
		Key    string `yaml:"key"`
		EndKey string `yaml:"endKey"`
	}
)

func (p *FixedConfig) GetKey() string {
	return p.Key
}

func (p *FixedConfig) OpOptions() clientv3.OpOption {
	return clientv3.WithLimit(0)
}

func (p *PrefixConfig) GetKey() string {
	return p.Key
}
func (p *PrefixConfig) OpOptions() clientv3.OpOption {
	return clientv3.WithPrefix()
}
func (p *RangeConfig) GetKey() string {
	return p.Key
}
func (p *RangeConfig) OpOptions() clientv3.OpOption {
	return clientv3.WithRange(p.EndKey)
}
