package agent

import (
	"fmt"
	"galaxy-guardian/logger"
	"galaxy-guardian/util"
	"galaxy-guardian/zookeeper"
)

func (a *Agent) GetVersion() {
	globalVersion, err := zookeeper.GetZKData(a.Config.AgentGlobalZK.Path, a.Config.AgentGlobalZK.Pair)
	if err != nil {
		logger.Logger.Errorf("get agent version from %s err: %v", a.Config.AgentGlobalZK.Path, err)
		return
	}
	a.globalVersion = string(globalVersion)

	targetVersion, _ := zookeeper.GetZKData(a.Config.AgentTargetZK.Path, a.Config.AgentTargetZK.Pair)
	a.targetVersion = string(targetVersion)
}

func (a *Agent) DownloadVersion(version string) {
	downloadUrl := fmt.Sprintf("%s/%s", a.Config.DownloadUrl, version)
	var downloadBinUrl, downloadAscUrl string
	var binFileName, ascFileName string
	if a.Config.Compress != "" {
		downloadBinUrl = fmt.Sprintf("%s.%s", downloadUrl, a.Config.Compress)
		binFileName = fmt.Sprintf("%s,%s", version, a.Config.Compress)
	} else {
		downloadBinUrl = fmt.Sprintf("%s.%s", downloadUrl, a.Config.Name)
		binFileName = a.Config.Name
	}

	cache, err := util.DownloadFromFileserver(downloadBinUrl, version, binFileName, a.Config.Name)
	if err != nil {
		logger.Logger.Errorf("download err: %v", err)
		return
	}

	if a.Config.RsaCheck {
		downloadAscUrl = fmt.Sprintf("%s.%s", downloadBinUrl, "asc")
		ascFileName = fmt.Sprintf("%s.%s", binFileName, "asc")
		ascCache, err := util.DownloadFromFileserver(downloadAscUrl, version, ascFileName, a.Config.Name)
		if err != nil {
			logger.Logger.Errorf("download err: %v", err)
			return
		}
		err = util.RsaSignVerify(a.Config.PublicKeyUrl, cache, ascCache, "agent")
	}

}
