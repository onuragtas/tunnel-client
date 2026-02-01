package tunnel

import (
	"fmt"
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
	_ = response
}

func login() {
	var username string
	var password string

	fmt.Println("Username:")
	fmt.Scanf("%s", &username)
	fmt.Println("Password:")
	fmt.Scanf("%s", &password)

	response := requestClient.Login(username, password)
	_ = response
}
