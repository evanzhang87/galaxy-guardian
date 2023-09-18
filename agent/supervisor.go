package agent

import (
	"errors"
	"fmt"
	"galaxy-guardian/logger"
	"galaxy-guardian/zookeeper"
	"github.com/samuel/go-zookeeper/zk"
	"sync"
)

type daemonType struct {
	Mutex *sync.Mutex
	Conn  *zk.Conn
}

var daemonQueue map[string]daemonType

func RegistryAgent(name, zkPair string) (error, *sync.Mutex) {
	if _, ok := daemonQueue[name]; !ok {
		conn := zookeeper.GetZkConn(zkPair)
		if conn == nil {
			return errors.New(fmt.Sprintf("can't get zk pair: %v", zkPair)), nil
		}
		mtx := &sync.Mutex{}
		daemonQueue[name] = daemonType{
			Mutex: mtx,
			Conn:  conn,
		}
		logger.Logger.Infof("add daemon: %v", name)
	}
	return nil, daemonQueue[name].Mutex
}
