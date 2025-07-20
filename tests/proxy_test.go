package tests

import (
	"testing"

	"github.com/onuragtas/tunnel-client"
	"github.com/onuragtas/tunnel-client/models"
)

var client = tunnel.NewClient()

func TestProxy(t *testing.T) {
	client.ListDomain()
	var c chan bool
	var tunnels []models.Tunnel

	tunnels = append(tunnels, models.Tunnel{
		LocalPort:     3333,
		DestinationIp: "127.0.0.1",
		DomainId:      409,
		LocalIp:       "127.0.0.1",
		Domain:        "test.tnpx.org",
	})

	client.StartTunnel(tunnels, "user", "pass")
	<-c
}

func TestAddDomain(t *testing.T) {
	client.ListDomain()
	client.CreateDomain("test")
}

func TestLogin(t *testing.T) {
	response := client.Login("user", "pass")
	t.Error(!response.Success)
}
