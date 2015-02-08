package plugin

import (
	plugin_clock "plugin/clock"
	"plugin_interface"
	"fmt"
)
func NewPlugin(
	name string, // Plugin name
	instance string, // Plugin instance
	config interface{}, // Plugin config
	update_filtered chan plugin_interface.Update, // Update channel
) (	chan plugin_interface.Event)  {
	events := make(chan plugin_interface.Event, 16)
	update := make(chan plugin_interface.Update,1)

	switch {
	case name == `clock`:
		go plugin_clock.New(config, events, update)
	case true:
		panic(fmt.Sprintf("no plugin named %s", name))
	}


	go filterUpdate(name, instance, update ,update_filtered)
	return events
}



func filterUpdate(
	name string,
	instance string,
	update chan plugin_interface.Update,
	update_filtered chan plugin_interface.Update ) {
	for {
		ev := <- update
		ev.Name = name
		ev.Instance = instance
		update_filtered <- ev
	}
}
