package main

import (
	"fmt"
	"tunnel-client/utils"
)

func register() {
	var username string
	var email string
	var password string

	fmt.Println("Email:")
	fmt.Scanf("%s", &email)
	fmt.Println("Username:")
	fmt.Scanf("%s", &username)
	fmt.Println("Password:")
	fmt.Scanf("%s", &password)

	response := requestClient.Register(username, password, email)
	utils.WriteToken(response.Data.Token)
}

func login() {
	var username string
	var password string

	fmt.Println("Username:")
	fmt.Scanf("%s", &username)
	fmt.Println("Password:")
	fmt.Scanf("%s", &password)

	response := requestClient.Login(username, password)
	utils.WriteToken(response.Data.Token)
}
