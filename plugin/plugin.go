package plugin

import (
	"github.com/meshnet-gophers/firmware/router"
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
)

type Plugin interface {
	Init(router *router.MeshRouter, u *pb.User, nodeID uint32)
	Name() string
}
