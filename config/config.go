package config

import (
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"os"
	"path/filepath"
)

var log = logging.MustGetLogger("main")

var cfgFiles =[]string{
	"$HOME/.config/uberstatus/uberstatus.conf",
	"./cfg/uberstatus.conf",
	"./cfg/uberstatus.default.conf",
	"/usr/share/doc/uberstatus/uberstatus.example.conf",
}



type PluginConfig struct {
	Name string
	Instance string
	Plugin string
	Config map[string]interface{}
}


type Config struct {
	Plugins []PluginConfig
}


func LoadConfig(file ...string) Config {
	var cfg Config
	var cfgFile string
	var cfgFileList []string
	if len(file) < 1 {
		cfgFileList = cfgFiles
	} else {
		cfgFileList = file
	}
	for _,element := range cfgFileList {
		filename, _ := filepath.Abs(os.ExpandEnv(element))
		log.Warning(filename)
		if _, err := os.Stat(filename); err == nil {
			cfgFile = filename
			break
		}
	}
	if cfgFile == "" {
		log.Panicf("could not find config file: %v", cfgFiles)
	}
	log.Infof("Loading config file: %s", cfgFile)
	raw_cfg, err := ioutil.ReadFile(cfgFile)
	err = yaml.Unmarshal([]byte(raw_cfg), &cfg)
	_ = err
	str, _ := yaml.Marshal(cfg)
	log.Debug(string(str))

	return cfg
}
