package tests

import (
	"github.com/onuragtas/tunnel-client"
	"github.com/onuragtas/tunnel-client/models"
	"testing"
)

var client = tunnel.NewClient()

func TestProxy(t *testing.T) {
	client.ListDomain()
	var c chan bool
	var tunnels []models.Tunnel

	tunnels = append(tunnels, models.Tunnel{
		IndexId:       22,
		LocalPort:     5000,
		DestinationIp: "ip",
		DomainId:      22,
		LocalIp:       "ip",
		Domain:        "testdomain",
	})

	client.StartTunnel(tunnels, "user", "pass")
	<-c
}

func TestAddDomain(t *testing.T) {
	client.ListDomain()
	client.CreateDomain("test")
}
