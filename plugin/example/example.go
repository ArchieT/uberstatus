package example


import (
	"github.com/XANi/uberstatus/uber"
//	"gopkg.in/yaml.v1"
	"time"
	"github.com/op/go-logging"
	"fmt"
)

// Example plugin for uberstatus
// plugins are wrapped in go() when loading

var log = logging.MustGetLogger("main")


// set up a config struct
type config struct {
	prefix string
	interval int
}

type state struct {
	cfg config
	cnt int
	ev int
}


func Run(c map[string]interface{}, events chan uber.Event, update chan uber.Update) {
	var st state
	st.cfg = loadConfig(c)
	// initial update on start
	update <- st.updatePeriodic()
	for {
		select {
		case updateEvent := (<-events):
			update <- st.updateFromEvent(updateEvent)
		case <-time.After(time.Duration(st.cfg.interval) * time.Millisecond):
			update <- st.updatePeriodic()
		}
	}
}


func (state *state) updatePeriodic() uber.Update {
	var update uber.Update
	update.FullText = fmt.Sprintf("nothing %d %d", state.cnt, state.ev)
	update.ShortText = `nope`
	update.Color = `#66cc66`
	state.cnt++
	return update
}

func (state *state) updateFromEvent(e uber.Event) uber.Update {
	var update uber.Update
	update.FullText = fmt.Sprintf("event: %+v", e)
	update.ShortText = `upd`
	update.Color = `#cccc66`
	state.ev++
	return update
}


// parse received structure into config
func loadConfig(c map[string]interface{}) config {
	var cfg config
	cfg.interval = 10000
	cfg.prefix = "ex: "
	for key, value := range c {
		converted, ok := value.(string)
		if ok {
			switch {
			case key == `prefix`:
				cfg.prefix=converted
			default:
				log.Warning("unknown config key: [%s]", key)

			}
		} else {
			converted, ok := value.(int)
			if ok {
				switch {
				case key == `interval`:
					cfg.interval = converted
				default:
					log.Warning("unknown config key: [%s]", key)
				}
			} else {
				log.Error("Cant interpret value of config key [%s]", key)
			}
		}
	}
	return cfg
}