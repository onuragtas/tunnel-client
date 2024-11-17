package tunnel

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type Endpoint struct {
	Host string
	Port int
}

type Client struct {
	LocalEndpoint     Endpoint
	ServerEndpoint    Endpoint
	RemoteEndpoint    Endpoint
	Signal            chan int
	CloseHandleSignal chan string
	SshUser           string
	SshPassword       string
}

func (t *Client) checkNetwork() (ok bool) {
	_, err := http.Get("http://clients3.google.com/generate_204")
	if err != nil {
		return false
	}
	return true
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

func (t *Client) Connect() {

	if !t.checkNetwork() {
		return
	}
	log.Println("starting...", t.RemoteEndpoint.Host)
	sshConfig := &ssh.ClientConfig{
		// SSH connection username
		User: t.SshUser,
		Auth: []ssh.AuthMethod{
			// put here your private key path
			ssh.Password(t.SshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", t.ServerEndpoint.String(), sshConfig)
	if err != nil {
		log.Println(fmt.Printf("Dial INTO remote server error: %s - %s ", err, t.ServerEndpoint.String()))
	}

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", t.RemoteEndpoint.String())
	if err != nil {
		time.Sleep(time.Second * 2)
		t.CloseHandleSignal <- t.RemoteEndpoint.Host
		log.Println(fmt.Printf("Listen open port ON remote server error: %s - %s ", err, t.RemoteEndpoint.String()))
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	go func(listener net.Listener) {
		<-t.Signal
		listener.Close()
		listener = nil
	}(listener)

	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", t.LocalEndpoint.String())
		if err != nil {
			log.Println(fmt.Printf("Dial INTO local service error: %s - %s ", err, t.LocalEndpoint.String()))
		}

		if listener == nil {
			break
		}
		client, err := listener.Accept()
		if err != nil {
			time.Sleep(time.Second * 2)
			t.CloseHandleSignal <- t.RemoteEndpoint.Host
			log.Println("kapandı", t.LocalEndpoint.String(), t.RemoteEndpoint.String(), t.ServerEndpoint.String())
			log.Println(err)
			break
		}

		go t.handleClient(client, local)
	}
}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func (t *Client) handleClient(client net.Conn, remote net.Conn) {

	defer func(client net.Conn) {
		if client != nil {
			err := client.Close()
			if err != nil {

			}
		}
	}(client)

	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		if !(remote == nil || client == nil) {
			_, err := io.Copy(client, remote)
			if err != nil {
				log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
			}
			fmt.Println("\rwhile copy remote->local:", client.RemoteAddr().String(), " -> ", client.LocalAddr().String(), " -> ", remote.LocalAddr().String(), " -> ", remote.RemoteAddr().String(), "\t\t\t\t\t\t\t\t\t")
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		if !(remote == nil || client == nil) {
			_, err := io.Copy(remote, client)
			if err != nil {
				log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
			}
			fmt.Println("\rwhile copy local->remote:", remote.RemoteAddr().String(), " -> ", remote.LocalAddr().String(), " -> ", client.LocalAddr().String(), " -> ", client.RemoteAddr().String(), "\t\t\t\t\t\t\t\t\t")
		}
		chDone <- true
	}()

	<-chDone
}
