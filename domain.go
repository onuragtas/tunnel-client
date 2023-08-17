package main

import (
	"fmt"
	"github.com/onuragtas/tunnel-client/models"
	tunnel2 "github.com/onuragtas/tunnel-client/tunnel"
	"github.com/onuragtas/tunnel-client/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Signal        chan string
	Domain        models.DomainItem
	KeepAliveTime time.Time
	Tunnel        tunnel2.Client
	CloseSignal   chan int
}

var closeHandleSignal = make(chan string)

func listDomain() {
	log.Println("Listing Domains...")
	response := requestClient.ListDomains(utils.ReadToken())
	for key, item := range response.Data.Domains {
		fmt.Println(key+1, item.Domain)
	}
	domainList = response
}

func listenClose() {
	for true {
		select {
		case domain := <-closeHandleSignal:
			for _, value := range startedTunnels.Data {
				if value.Domain.Domain == domain {
					go value.Tunnel.Connect()
				}
			}
		}
	}
}

func createDomain() {
	fmt.Println(requestClient.CreateNewDomain(utils.ReadToken()))
}

func startTunnel() {
	listDomain()
	log.Println("Select IDs, Ex: 1:local_port:destination_ip,2:local_port:destination_ip...Or 1,2...\nClose:0")

	if list == "" {
		fmt.Scanf("%s", &list)
	}

	if list != "0" {
		ids := strings.Split(list, ",")
		for _, item := range ids {
			var tunnel tunnel2.Client
			splitted := strings.Split(item, ":")

			item = splitted[0]
			localPort := defaultLocalPort
			localIp := defaultDestinationIp
			if len(splitted) > 1 {
				localPort = splitted[1]
			}
			if len(splitted) > 2 {
				localIp = splitted[2]
			}

			c := make(chan int)

			tunnelDetail := getTunnelItem(item)

			if tunnelDetail != nil {
				continue
			}

			domainDetail := getDomain(item)

			// local service to be forwarded
			destinationLocalPort, _ := strconv.Atoi(localPort)
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
func renewDomain() {
	var id string
	listDomain()
	log.Println("Select ID, Ex: 1\nClose:0")

	if id == "" {
		fmt.Scanf("%s", &id)
	}

	if id != "0" {
		domainDetail := getDomain(id)

		requestClient.RenewDomain(utils.ReadToken(), domainDetail.Domain)
	}

}

func closeTunnel() {
	listDomain()
	log.Println("Select IDs, Ex: 1,2,3,4\nClose:0")

	if list == "" {
		fmt.Scanf("%s", &list)
	}

	if list != "0" {
		ids := strings.Split(list, ",")
		for _, item := range ids {
			tunnelDetail := getTunnelItem(item)
			if tunnelDetail != nil {
				tunnelDetail.CloseSignal <- 1
				removeTunnelItem(item)
			}
		}
	}
}

func deleteDomain() {
	listDomain()
	log.Println("Select IDs, Ex: 1,2,3,4\nClose: 0")

	var idList []string

	if list == "" {
		fmt.Scanf("%s", &list)
	}

	if list != "0" {
		ids := strings.Split(list, ",")
		for _, item := range ids {
			idList = append(idList, strconv.Itoa(getDomain(item).ID))
		}
		fmt.Println(requestClient.DeleteDomain(utils.ReadToken(), idList))
	}

}

func getDomain(id string) models.DomainItem {
	for key, item := range domainList.Data.Domains {
		idCasting, _ := strconv.Atoi(id)
		if key == idCasting-1 {
			return item
		}
	}
	return models.DomainItem{}
}

func getTunnelItem(id string) *Item {
	for key, item := range startedTunnels.Data {
		idCasting, _ := strconv.Atoi(id)
		if key == idCasting-1 {
			return &item
		}
	}
	return nil
}

func removeTunnelItem(id string) *Item {
	for key, _ := range startedTunnels.Data {
		idCasting, _ := strconv.Atoi(id)
		if key == idCasting-1 {
			startedTunnels.Data = remove(startedTunnels.Data, key)
		}
	}
	return nil
}

func remove(slice []Item, s int) []Item {
	return append(slice[:s], slice[s+1:]...)
}
