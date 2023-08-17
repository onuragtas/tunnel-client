package tunnel

import (
	"fmt"
	"github.com/onuragtas/tunnel-client/models"
	"github.com/onuragtas/tunnel-client/utils"
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

func (c *Client) StartTunnel() {

}

func (c *Client) DeleteDomain(idList []string) models.Response {
	return requestClient.DeleteDomain(c.GetToken(), idList)
}

func (c *Client) CloseDomain() {

}

func (c *Client) RenewDomain(domain string) {
	requestClient.RenewDomain(utils.ReadToken(), domain)
}
