package zookeeper

import (
	"errors"
	"galaxy-guardian/logger"
	"github.com/samuel/go-zookeeper/zk"
	"time"
)

var connMap = make(map[string]*zk.Conn)

func InitZk(name string, addr []string) error {
	conn, events, err := zk.Connect(addr, time.Second*5)
	if err != nil {
		return err
	}
	ticker := time.NewTimer(time.Second * 15)
	defer ticker.Stop()
	for {
		select {
		case event := <-events:
			logger.Logger.Infof("get zk state: %v", event.State)
			if event.State == zk.StateHasSession {
				connMap[name] = conn
				return nil
			}
		case <-ticker.C:
			logger.Logger.Errorf("connect to zk %v timeout !", addr)
			return errors.New("zk timeout")
		}
	}
}

func GetZkConn(name string) *zk.Conn {
	return connMap[name]
}

func GetZKData(name, path string) (value []byte, err error) {
	if conn, ok := connMap[name]; ok {
		value, _, err = conn.Get(path)
		return
	}
	return
}
