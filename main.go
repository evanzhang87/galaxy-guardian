package main

import (
	"encoding/json"
	"flag"
	"fmt"
	agent2 "galaxy-guardian/agent"
	conf "galaxy-guardian/config"
	"galaxy-guardian/logger"
	"galaxy-guardian/quota"
	"galaxy-guardian/zookeeper"
	"log"
	"os"
	"reflect"
	"runtime"
	"sync"
)

var (
	config  string
	version bool
)

func init() {
	flag.StringVar(&config, "c", "config.yaml", "config file path")
	flag.BoolVar(&version, "v", false, "print version")
}

func main() {
	flag.Parse()
	if version {
		fmt.Printf("version: %s\ngo version: %s\nbuild time:%s\n", Version, runtime.Version(), BuildTime)
		return
	}
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(1)
	}
	conf.InitConfig(config)
	Run()
}

func Run() {
	c := conf.GetConfig()
	if c == nil {
		log.Fatalf("can't load core config")
	}
	go quota.WatchRss(c.RssLimitMb)
	go quota.WatchSocket(c.SocketLimit)

	for _, pair := range c.ZkPairs {
		err := zookeeper.InitZk(pair.Name, pair.Hosts)
		if err != nil {
			logger.Logger.Errorf("init zk %v err: %v", pair.Name, err)
		}
	}

	confExpand := make(map[string]string)
	confExpand["registry_url"] = c.RegistryUrl
	confExpand["fileserver_url"] = c.FileserverUrl
	for _, agentConf := range c.Agents {
		dirtyJsonByte, _ := json.Marshal(agentConf)
		var temp map[string]interface{}
		_ = json.Unmarshal(dirtyJsonByte, &temp)
		for k, v := range temp {
			if v == nil {
				continue
			}
			if reflect.TypeOf(v).Kind() == reflect.String {
				confExpand[k] = v.(string)
			} else if reflect.TypeOf(v).Kind() == reflect.Map {
				innerTemp := v.(map[string]interface{})
				for ik, iv := range innerTemp {
					if reflect.TypeOf(iv).Kind() == reflect.String {
						confExpand[fmt.Sprintf("%s.%s", k, ik)] = iv.(string)
					}
				}
			}
		}
		err := agentConf.FixConfig(confExpand)
		if err != nil {
			log.Fatalf("init config error: %v", err)
		}
		agent := agent2.Agent{
			Mutex:  &sync.Mutex{},
			Config: agentConf,
		}
		go agent.Run()
	}
}
