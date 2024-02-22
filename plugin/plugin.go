package plugin

var plugins []Plugin

func Register(p Plugin) {

}

type Plugin interface {
	Init()
}
