// NOT true random number generation on RP2040 but good enough for quick demos.
// Do not use in important production systems.
// Seriously, don't.
package insecure_rng

import (
	"github.com/meshnet-gophers/firmware/plugin/registry"
	"github.com/meshnet-gophers/firmware/router"
	pb "github.com/meshnet-gophers/meshtastic-go/meshtastic"
	"machine"

	"crypto/rand"
)

func init() {
	r := reader{}
	rand.Reader = &r
	registry.RegisterPlugin(&r)
}

func (b *reader) Name() string {
	return "insecure_rng"
}

func (b *reader) Init(r *router.MeshRouter, u *pb.User, nodeID uint32) {}

type reader struct{}

func (r *reader) Read(b []byte) (n int, err error) {
	if len(b) == 0 {
		return
	}

	var randomByte uint32
	for i := range b {
		if i%4 == 0 {
			randomByte, err = machine.GetRNG()
			if err != nil {
				return n, err
			}
		} else {
			randomByte >>= 8
		}
		b[i] = byte(randomByte)
	}

	return len(b), nil
}
