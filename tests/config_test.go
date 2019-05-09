package tests

import (
	"fmt"
	"github.com/blademainer/etcd-sync/etcdsync"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestMarshal(t *testing.T) {

	syncConfig := &etcdsync.MasterConfig{
		FixedConfigs: []*etcdsync.FixedConfig{
			{"/com"},
		},
		PrefixConfigs: []*etcdsync.PrefixConfig{{""}},
		RangeConfigs: []*etcdsync.RangeConfig{
			{"/com/a", "/com/b"},
		},
	}
	syncConfig.Endpoints = []string{"http://10.10.134.30:2379", "http://10.10.134.31:2379", "http://10.10.134.32:2379"}
	syncConfig.TimeoutSeconds = 10

	config := &etcdsync.EtcdSync{
		Master: syncConfig,
		Slaves: []*etcdsync.EtcdConfig{
			{
				Endpoints:      []string{"http://10.10.134.30:2379", "http://10.10.134.31:2379", "http://10.10.134.32:2379"},
				TimeoutSeconds: 10,
			},
			{
				Endpoints:      []string{"http://10.10.134.30:12379", "http://10.10.134.31:12379", "http://10.10.134.32:12379"},
				TimeoutSeconds: 10,
			},
		},
	}

	out, _ := yaml.Marshal(config)
	expect := string(out)
	fmt.Println("################################ config ################################")
	fmt.Println(expect)
	fmt.Println("########################################################################")

	sync := &etcdsync.EtcdSync{}
	yaml.Unmarshal(out, sync)
	bytes, _ := yaml.Marshal(config)
	result := string(bytes)
	if expect != result {
		t.Fail()
	}

}
