package config

import (
	"errors"
	"fmt"
	"galaxy-guardian/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var Config *CoreConfig
var pathRegex = regexp.MustCompile(`\{([^}]+)\}`)

type CoreConfig struct {
	RegistryUrl   string         `yaml:"registry_url"`
	FileserverUrl string         `yaml:"fileserver_url"`
	PprofPort     int            `yaml:"pprof_port"`
	RssLimitMb    int            `yaml:"rss_limit_mb"` // rss_limit_mb to kill client, default is 100Mb
	SocketLimit   int            `yaml:"socket_limit"` // 100,just set in host config,not set in ClientConfig
	LogLevel      string         `yaml:"log_level"`
	LogFile       string         `yaml:"log_file"`
	LogMaxsize    int            `yaml:"log_maxsize"` // size in kb
	LogMaxrolls   int            `yaml:"log_maxrolls"`
	ZkPairs       []*ZkPair      `yaml:"zk_pairs"`
	Agents        []*AgentConfig `yaml:"agents"`
}

type ZkPair struct {
	Name  string   `yaml:"name"`
	Hosts []string `yaml:"hosts"`
}

type AgentConfig struct {
	Name           string          `yaml:"name" json:"name"`         // name of bin
	Source         string          `yaml:"source" json:"source"`     // bin|file|yum|apt
	Compress       string          `yaml:"compress" json:"compress"` // tar.gz
	Workdir        string          `yaml:"workdir" json:"workdir"`
	IsDaemon       bool            `yaml:"is_daemon" json:"is_daemon"`
	InstallCmd     string          `yaml:"install_cmd" json:"install_cmd"`
	RunCmd         string          `yaml:"run_cmd" json:"run_cmd"` // eg: {workdir}/{name} service start | systemctl start {name}
	StopCmd        string          `yaml:"stop_cmd" json:"stop_cmd"`
	RestartCmd     string          `yaml:"restart_cmd" json:"restart_cmd"`
	AgentTargetZK  *ZkUnit         `yaml:"agent_target_zk" json:"agent_target_zk"`
	AgentGlobalZK  *ZkUnit         `yaml:"agent_global_zk" json:"agent_global_zk"`
	DownloadUrl    string          `yaml:"download_url" json:"download_url"` // url for download agent
	PublicKeyUrl   string          `yaml:"public_key_url" json:"public_key_url"`
	HeartbeatUrl   string          `yaml:"heartbeat_url" json:"heartbeat_url"`
	RsaCheck       bool            `yaml:"rsa_check" json:"rsa_check"` // check download file
	ZkPairName     string          `yaml:"zk_pair_name" json:"zk_pair_name"`
	HashCheck      bool            `yaml:"hash_check" json:"hash_check"` // check config hash
	AgentBinConfig *AgentBinConfig `yaml:"agent_bin_config" json:"agent_bin_config"`
}

type AgentBinConfig struct {
	ConfigDir      string  `yaml:"config_dir" json:"config_dir"`
	ConfigZkPair   string  `yaml:"config_zk_pair" json:"config_zk_pair"`
	ConfigTargetZK *ZkUnit `yaml:"config_target_zk" json:"config_target_zk"`
	ConfigGlobalZK *ZkUnit `yaml:"config_global_zk" json:"config_global_zk"`
}

func (a *AgentConfig) FixConfig(confMap map[string]string) error {
	fixedRunCmd, err := fixSchema(a.RunCmd, confMap)
	if err != nil {
		return err
	}
	a.RunCmd = fixedRunCmd
	fixedStopCmd, err := fixSchema(a.StopCmd, confMap)
	if err != nil {
		return err
	}
	a.StopCmd = fixedStopCmd
	fixedRestart, err := fixSchema(a.RestartCmd, confMap)
	if err != nil {
		return err
	}
	a.RestartCmd = fixedRestart
	fixedDownloadUrl, err := fixSchema(a.DownloadUrl, confMap)
	if err != nil {
		return err
	}
	a.DownloadUrl = fixedDownloadUrl
	fixedHeartbeatUrl, err := fixSchema(a.HeartbeatUrl, confMap)
	if err != nil {
		return err
	}
	a.HeartbeatUrl = fixedHeartbeatUrl
	return nil
}

func fixSchema(value string, confMap map[string]string) (string, error) {
	if strings.Contains(value, "{") && strings.Contains(value, "}") {
		res := pathRegex.FindAllString(value, -1)
		for _, r := range res {
			if v, ok := confMap[r[1:len(r)-1]]; ok {
				value = strings.Replace(value, r, v, 1)
			}
		}
	}
	if strings.Contains(value, "{") && strings.Contains(value, "}") {
		return "", errors.New(fmt.Sprintf("can't merge config: %v", value))
	}
	return value, nil
}

type ZkUnit struct {
	Path string `yaml:"path" json:"path"`
	Pair string `yaml:"pair" json:"pair"`
}

func InitConfig(config string) {
	content, _ := ioutil.ReadFile(config)
	if len(content) != 0 {
		err := yaml.Unmarshal(content, Config)
		if err != nil {
			log.Fatalf("can't load config from %v, exit!", Config)
		}
	}
	Config.fix()
	logger.Logger.Info("init logger...")
	logger.InitLogger(Config.LogMaxsize, Config.LogMaxrolls, Config.LogFile, Config.LogLevel)
}

func LoadConfig(config *CoreConfig) {
	Config = config
	Config.fix()
	logger.Logger.Info("init logger...")
	logger.InitLogger(Config.LogMaxsize, Config.LogMaxrolls, Config.LogFile, Config.LogLevel)
}

func GetConfig() *CoreConfig {
	return Config
}

func (c *CoreConfig) fix() {
	if len(c.ZkPairs) == 0 {
		log.Fatal("can't load zk config, exit!")
	}
	if c.RssLimitMb == 0 {
		c.RssLimitMb = 100
	}
	if c.SocketLimit == 0 {
		c.SocketLimit = 100
	}
	if c.LogLevel == "" {
		c.LogLevel = "INFO"
	}
	if c.LogFile == "" {
		c.LogFile = "guardian.log"
	}
	if c.LogMaxsize == 0 {
		c.LogMaxsize = 10240
	}
	if c.LogMaxrolls == 0 {
		c.LogMaxrolls = 7
	}
}
