package agent

import (
	"fmt"
	"galaxy-guardian/config"
	"galaxy-guardian/logger"
	"galaxy-guardian/util"
	"os"
	"sync"
	"time"
)

const defaultTimeout = time.Second * 5

type Agent struct {
	Mutex         *sync.Mutex
	Config        *config.AgentConfig
	pid           string
	globalVersion string
	targetVersion string
}

func (a *Agent) Run() {
	err, _ := RegistryAgent(a.Config.Name, a.Config.ZkPairName)
	if err != nil {
		logger.Logger.Errorf("init agent err: %v", err)
		return
	}
	a.CheckStatus()
	a.GetVersion()
	if a.pid == "" {
		logger.Logger.Infof("agent %s is not running, try to start", a.Config.Name)
		_, err = os.Stat(fmt.Sprintf("%s/%s", a.Config.Workdir, a.Config.Name))
		if err == nil {
			logger.Logger.Infof("agent %s bin is exist", a.Config.Name)
		} else if os.IsNotExist(err) {
			logger.Logger.Infof("agent %s need download", a.Config.Name)
			if a.targetVersion != "" {
				a.DownloadVersion(a.targetVersion)
			} else {
				a.DownloadVersion(a.globalVersion)
			}
		} else {
			logger.Logger.Errorf("get agent bin file err: %v", err)
			return
		}

		if a.Config.InstallCmd != "" {

		}
	}
}

func (a *Agent) CheckStatus() {
	proFilter := fmt.Sprintf("grep %s/%s", a.Config.Workdir, a.Config.Name)
	if a.Config.AgentBinConfig.ConfigDir != "" {
		proFilter = proFilter + " | grep " + a.Config.AgentBinConfig.ConfigDir
	}
	pscmd := fmt.Sprintf("ps -ef | grep %s | awk '{print $2}'", proFilter)
	pidByte, err := util.Command(defaultTimeout, pscmd)
	if err != nil {
		logger.Logger.Errorf("get pid err: %v", err)
		return
	}
	a.pid = string(pidByte)
	logger.Logger.Infof("get running agent: %s [%s]", a.Config.Name, a.pid)
}

func (a *Agent) Install() {

}

func (a *Agent) Start() {

}

func (a *Agent) Stop() {

}

func (a *Agent) Restart() {

}
