package plugin






import (
	"fmt"
	"github.com/op/go-logging"
	"gopkg.in/yaml.v1"
)
var log = logging.MustGetLogger("main")

type Plugin interface {
//	New(config map[string]interface{}, update chan plugin.Update) error
	SendEvent(*plugin.Event) error
	Start() (chan Event)
}
func NewPlugin(f func(string, interface{}) Plugin, addr string, cfg interface{}) Plugin {
	t := f(addr, cfg)
	return t
}
