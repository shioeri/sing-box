package main

import (
	"net/netip"
	"testing"
	"time"

	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common"
	"github.com/sagernet/sing/common/json/badoption"
)

func _TestWireGuard(t *testing.T) {
	startDockerContainer(t, DockerOptions{
		Image: ImageBoringTun,
		Cap:   []string{"MKNOD", "NET_ADMIN", "NET_RAW"},
		Ports: []uint16{serverPort, testPort},
		Bind: map[string]string{
			"wireguard.conf": "/etc/wireguard/wg0.conf",
		},
		Cmd: []string{"wg0"},
	})
	time.Sleep(5 * time.Second)
	startInstance(t, option.Options{
		Inbounds: []option.Inbound{
			{
				Type: C.TypeMixed,
				Options: &option.HTTPMixedInboundOptions{
					ListenOptions: option.ListenOptions{
						Listen:     common.Ptr(badoption.Addr(netip.IPv4Unspecified())),
						ListenPort: clientPort,
					},
				},
			},
		},
		Outbounds: []option.Outbound{
			{
				Type: C.TypeWireGuard,
				Options: &option.WireGuardEndpointOptions{
					ServerOptions: option.ServerOptions{
						Server:     "127.0.0.1",
						ServerPort: serverPort,
					},
					Address:       []netip.Prefix{netip.MustParsePrefix("10.0.0.2/32")},
					PrivateKey:    "qGnwlkZljMxeECW8fbwAWdvgntnbK7B8UmMFl3zM0mk=",
					PeerPublicKey: "QsdcBm+oJw2oNv0cIFXLIq1E850lgTBonup4qnKEQBg=",
				},
			},
		},
	})
	testSuitWg(t, clientPort, testPort)
}
