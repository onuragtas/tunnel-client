package tunnel

import (
	"fmt"
	"github.com/onuragtas/tunnel-client/models"
	tunnel2 "github.com/onuragtas/tunnel-client/tunnel"
	"github.com/onuragtas/tunnel-client/utils"
	"time"
)

type IClient interface {
	Login() error
	Register() error
	CheckUser() bool
	GetToken() string
	ListDomain() models.Domain
	CreateDomain()
	StartTunnel()
	CloseTunnel()
	DeleteDomain()
	RenewDomain()
}

type Client struct {
	IClient
}

func NewClient() *Client {
	return &Client{}
}

func (c Client) GetToken() string {
	return utils.ReadToken()
}

func (c Client) initialize() {
	if utils.ReadToken() == "" {
		c.CheckUser()
	}
}

func (c *Client) Login(username, password string) bool {
	response := requestClient.Login(username, password)
	utils.WriteToken(response.Data.Token)
	return true
}

func (c *Client) Register(username, password, email string) bool {
	response := requestClient.Register(username, password, email)
	utils.WriteToken(response.Data.Token)
	return true
}

func (c *Client) CheckUser() bool {
	if utils.ReadToken() == "" {
		return false
	}

	return true
}

func (c *Client) ListDomain() models.Domain {
	response := requestClient.ListDomains(utils.ReadToken())
	for key, item := range response.Data.Domains {
		fmt.Println(key+1, item.Domain)
	}
	domainList = response
	return domainList
}

func (c *Client) CreateDomain(domain string) interface{} {
	return requestClient.CreateNewDomain(domain, utils.ReadToken())
}

func (c *Client) StartTunnel(tunnelList []models.Tunnel, sshUser, sshPassword string) {

	if len(tunnelList) != 0 {
		for _, item := range tunnelList {
			var tunnel tunnel2.Client

			localPort := defaultLocalPort
			localIp := defaultDestinationIp
			if item.LocalPort != 0 {
				localPort = item.LocalPort
			}

			if item.LocalIp != "" {
				localIp = item.LocalIp
			}

			c := make(chan int)

			tunnelDetail := getTunnelItem(item.IndexId)

			if tunnelDetail != nil {
				continue
			}

			domainDetail := getDomain(item.DomainId)

			// local service to be forwarded
			destinationLocalPort := localPort
			var localEndpoint = tunnel2.Endpoint{
				Host: localIp,
				Port: destinationLocalPort,
			}

			// remote SSH server
			var serverEndpoint = tunnel2.Endpoint{
				Host: domainDetail.Domain,
				Port: 22,
			}

			// remote forwarding port (on remote SSH server network)
			var remoteEndpoint = tunnel2.Endpoint{
				Host: domainDetail.Domain,
				Port: domainDetail.Port,
			}

			tunnel.LocalEndpoint = localEndpoint
			tunnel.RemoteEndpoint = remoteEndpoint
			tunnel.ServerEndpoint = serverEndpoint
			tunnel.Signal = c
			tunnel.CloseHandleSignal = closeHandleSignal
			tunnel.SshUser = sshUser
			tunnel.SshPassword = sshPassword

			startedTunnels.Data = append(startedTunnels.Data, Item{Signal: closeHandleSignal, CloseSignal: c, Domain: domainDetail, KeepAliveTime: time.Now(), Tunnel: tunnel})

			go tunnel.Connect()
		}
	}

	go listenClose()
}

func (c *Client) DeleteDomain(idList []string) models.Response {
	return requestClient.DeleteDomain(c.GetToken(), idList)
}

func (c *Client) CloseDomain(closeList []int) {
	if len(closeList) != 0 {
		for _, item := range closeList {
			tunnelDetail := getTunnelItem(item)
			if tunnelDetail != nil {
				tunnelDetail.CloseSignal <- 1
				removeTunnelItem(item)
			}
		}
	}
}

func (c *Client) RenewDomain(domain string) {
	requestClient.RenewDomain(utils.ReadToken(), domain)
}
