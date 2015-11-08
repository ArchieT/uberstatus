package main

import (
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
	//	"io/ioutil"
	"os"
	"regexp"
	"time"
	//
	"github.com/XANi/uberstatus/config"
	"github.com/XANi/uberstatus/i3bar"
	"github.com/XANi/uberstatus/plugin"
	"github.com/XANi/uberstatus/uber"
)

type Config struct {
	Plugins *map[string]map[string]interface{}
}

var log = logging.MustGetLogger("main")
var logFormat = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortpkg}↛%{shortfunc}: %{level:.4s} %{id:03x} %{color:reset}%{message}",
)

type pluginMap struct {
	// channels used to send events to plugin
	input   map[string]map[string]chan uber.Event
	output  map[string]map[string]chan uber.Update
	slots   []i3bar.Msg
	slotMap map[string]map[string]int
}

func main() {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logBackendFormatter := logging.NewBackendFormatter(logBackend, logFormat)
	_ = logBackendFormatter
	logBackendLeveled := logging.AddModuleLevel(logBackendFormatter)
	logBackendLeveled.SetLevel(logging.DEBUG, "")
	//	logBackendLeveled.SetLevel(logging.NOTICE, "")
	logging.SetBackend(logBackendLeveled)

	log.Info("Starting")
	header := i3bar.NewHeader()
	msg := i3bar.NewMsg()
	msg.FullText = `test`
	b, err := json.Marshal(header)
	if err != nil {
		fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
	//c, err := json.Marshal(msg)

	i3input := i3bar.EventReader()
	updates := make(chan uber.Update, 10)
	cfg := config.LoadConfig()
	plugins := pluginMap{
		slotMap: make(map[string]map[string]int),
		slots:   make([]i3bar.Msg, len(cfg.Plugins)),
		input:   make(map[string]map[string]chan uber.Event),
	}
	for idx, pluginCfg := range cfg.Plugins {
		log.Info("Loading plugin %s into slot %d: %+v", pluginCfg.Plugin, idx, pluginCfg)
		if plugins.slotMap[pluginCfg.Name] == nil {
			plugins.slotMap[pluginCfg.Name] = make(map[string]int)
			plugins.input[pluginCfg.Name] = make(map[string]chan uber.Event)
		}
		//		if len(pluginCfg.Instance) == 0 {
		//			pluginCfg.Instance = pluginCfg.Name
		//		}

		plugins.slotMap[pluginCfg.Name][pluginCfg.Instance] = idx
		plugins.slots[idx] = i3bar.NewMsg()
		plugins.input[pluginCfg.Name][pluginCfg.Instance] = plugin.NewPlugin(pluginCfg.Name, pluginCfg.Instance, pluginCfg.Plugin, pluginCfg.Config, updates)
	}

	// fmt.Println("\n[")

	// plugins := config.Plugins
	// ifd := (*plugins)[`clock`] //.(map[string]interface{})
	// net := (*plugins)[`clock`] //.(map[string]interface{})
	// //	_ = plugin.NewPlugin("clock", "", &ifd, updates)
	// _ = plugin.NewPlugin("network", "", &net, updates)
	// _ = ifd
	fmt.Println(`[`)

	for {
		select {
		case ev := (<-i3input):
			plugins.parseEvent(ev)
		case upd := <-updates:
			plugins.parseUpdate(upd)
		case <-time.After(time.Second * 10):
			log.Info("Time passed")
		}
		fmt.Print(`[`)

		for idx, msg := range plugins.slots {
			os.Stdout.Write(msg.Encode())
			if idx+1 < (len(plugins.slots)) {
				fmt.Print(`,`)
			}
		}
		fmt.Println(`],`)

	}
}

func (plugins *pluginMap) parseUpdate(update uber.Update) {
	if val, ok := plugins.slotMap[update.Name][update.Instance]; ok {
		plugins.slots[val] = i3bar.CreateMsg(update)
	} else {
		log.Warning("Got msg from unknown place, name: %s, instance: %s", update.Name, update.Instance)
	}
}

func (plugins *pluginMap) parseEvent(ev uber.Event) {
	if val, ok := plugins.slotMap[ev.Name][ev.Instance]; ok {
		log.Info("got event for %+v", val)
	} else {
		log.Info("rejected event %+v", ev)
	}

}

func getTime() []byte {
	msg := i3bar.NewMsg()
	msg.Name = "clock"
	t := time.Now().Local()
	// reference Mon Jan 2 15:04:05 MST 2006 (unix: 1136239445)
	msg.FullText = t.Format(`15:04:05`)
	msg.Color = `#ffffff`
	return msg.Encode()
}

func San(in []byte) []byte {
	re := regexp.MustCompile(`\,{`)
	return re.ReplaceAllLiteral(in, []byte(`{`))
}

func PrintInterface(a interface{}) {
	fmt.Println("Interface:")
	txt, _ := yaml.Marshal(a)
	fmt.Printf("%s", txt)
}
