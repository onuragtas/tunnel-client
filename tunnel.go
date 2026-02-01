package tunnel

import (
	"time"

	"github.com/onuragtas/tunnel-client/models"
	tunnel2 "github.com/onuragtas/tunnel-client/tunnel"
	"github.com/onuragtas/tunnel-client/utils"
)

type IClient interface {
	Login() error
	Register() error
	CheckUser(...*string) bool
	GetToken(*string) string
	ListDomain(*string) models.Domain
	CreateDomain(string, *string) interface{}
	StartTunnel()
	CloseTunnel()
	DeleteDomain([]string, *string) models.Response
	RenewDomain(string, *string)
}

type Client struct {
	IClient
}

func NewClient() *Client {
	return &Client{}
}

// getToken: parametre olarak token gelirse onu döner, nil/boş gelirse diskten okur.
func getToken(token *string) string {
	if token != nil && *token != "" {
		return *token
	}
	return utils.ReadToken()
}

func (c Client) GetToken(token *string) string {
	return getToken(token)
}

func (c *Client) Login(username, password string) models.Login {
	response := requestClient.Login(username, password)
	utils.WriteToken(response.Data.Token)
	return response
}

func (c *Client) Logout() bool {
	utils.WriteToken("")
	return true
}

func (c *Client) Register(username, password, email string) models.Register {
	response := requestClient.Register(username, password, email)
	utils.WriteToken(response.Data.Token)
	return response
}

func (c *Client) CheckUser(token *string) bool {
	if getToken(token) == "" {
		return false
	}

	return true
}

func (c *Client) ListDomain(token *string) models.Domain {
	response := requestClient.ListDomains(getToken(token))
	domainList = response
	return domainList
}

func (c *Client) UserInfo(token *string) models.UserInfo {
	return requestClient.UserInfo(getToken(token))
}

func (c *Client) CreateDomain(domain string, token *string) interface{} {
	return requestClient.CreateNewDomain(domain, getToken(token))
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

			tunnelDetail := getTunnelItem(item.Domain)

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
				Port: 8222,
			}

			// remote forwarding port (on remote SSH server network)
			var remoteEndpoint = tunnel2.Endpoint{
				Host: "127.0.0.1",
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

func (c *Client) DeleteDomain(idList []string, token *string) models.Response {
	return requestClient.DeleteDomain(getToken(token), idList)
}

func (c *Client) CloseTunnel(closeList []string) {
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

func (c *Client) RenewDomain(domain string, token *string) {
	requestClient.RenewDomain(getToken(token), domain)
}

func (c *Client) GetStartedTunnels() StartedTunnels {
	return startedTunnels
}
