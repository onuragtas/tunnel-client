package main

import (
	"fmt"
	"tunnel-client/utils"
)

func main() {

	for true {

		if utils.ReadToken() == "" {
			checkUser()
		} else {
			listUserProcesses()
		}

	}
}

func listUserProcesses() {

	fmt.Println("\nProcess:\nList Domains: 1\nCreate New: 2\nStart Tunnel: 3\nDelete Domain:4\nClose Tunnel: 5\nRenew Domain: 6\n\n")
	if process == "" {
		fmt.Print("Process Number:")
		fmt.Scanf("%s", &process)
	}
	switch process {
	case "1":
		listDomain()
		break
	case "2":
		createDomain()
		break
	case "3":
		startTunnel()
		break
	case "4":
		deleteDomain()
		break
	case "5":
		closeTunnel()
		break
	case "6":
		renewDomain()
		break
	default:
		listDomain()
		break
	}

	process = ""
}

func checkUser() {

	var process string

	fmt.Println("You must be logged in.")
	fmt.Println("Process:\nLogin: 1\nRegister: 2")
	fmt.Scanf("%s", &process)
	switch process {
	case "1":
		login()
		break
	case "2":
		register()
		break
	}

	process = ""
}
