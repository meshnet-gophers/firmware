package registry

import (
	"github.com/meshnet-gophers/firmware/plugin"
	"github.com/meshnet-gophers/firmware/router"
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
	"log/slog"
)

var plugins []plugin.Plugin

func RegisterPlugin(p plugin.Plugin) {
	plugins = append(plugins, p)
}

func StartPlugins(router *router.MeshRouter, u *pb.User, nodeID uint32) {
	for _, p := range plugins {
		slog.Info("starting plugin", "name", p.Name())
		p.Init(router, u, nodeID)
	}
}
