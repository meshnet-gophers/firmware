package beacon

import "github.com/meshnet-gophers/firmware/plugin"

func init() {
	plugin.Register(new(Beacon))
}

type Beacon struct {
}

func (b Beacon) Init() {
	//TODO implement me
	panic("implement me")
}
