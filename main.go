package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/blademainer/etcd-sync/etcdsync"
	"github.com/blademainer/etcd-sync/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
)

// Unmarshal decodes the first document found within the in byte slice
// and assigns decoded values into the out value.
func UnmarshalFromFile(name string, out interface{}) error {
	if out == nil {
		return errors.New("out is nil")
	}

	// Check if the file is existed
	if not_existed, err := IsNotExist(name); not_existed {
		return err
	}

	// Read the whole file
	stream, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	// Unmarshal the yaml stream
	return yaml.Unmarshal(stream, out)
}

// Check if the file or directory is existed
func IsNotExist(name string) (bool, error) {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return true, err
		}
	}

	return false, nil
}

func main() {
	go func() {
		// open: http://localhost:6060/debug/pprof/
		result := http.ListenAndServe("localhost:6060", nil)
		fmt.Println("init result: ", result)
	}()

	var configFile string
	flag.StringVar(&configFile, "c", "conf/etcd-sync.yml", "Config file location")
	flag.Parse()

	syncConfig := &etcdsync.EtcdSync{}
	if err := UnmarshalFromFile(configFile, syncConfig); err != nil {
		panic(err)
	}
	loggerConfig := &logger.LoggerConfig{}
	if err := UnmarshalFromFile("conf/logger.yml", loggerConfig); err != nil {
		panic(err)
	}

	logger.Log.Init(*loggerConfig)

	logger.Log.Infof("Read config： %v", syncConfig)
	syncConfig.Start()
	logger.Log.Errorf("Program exit!!!")
}
