package request

import (
	"encoding/json"
	"fmt"
	"github.com/onuragtas/go-requests"
	"tunnel-client/models"
)

type Request struct {
	BaseUrl string
}

func (t *Request) post(endpoint string, parameters map[string]interface{}) []byte {
	request := requests.Request{BaseUrl: t.BaseUrl + endpoint, Parameters: parameters, Headers: map[string]string{"Content-Type": "application/json"}}
	request.Post()
	return request.GetBody()
}

func (t *Request) Login(username, password string) models.Login {
	body := make(map[string]interface{})
	body["username"] = username
	body["password"] = password

	var res models.Login

	response := t.post("/login", body)
	json.Unmarshal(response, &res)
	return res
}

func (t *Request) Register(username string, password string, email string) models.Register {
	body := make(map[string]interface{})
	body["username"] = username
	body["password"] = password
	body["email"] = email
	var res models.Register

	response := t.post("/register", body)
	json.Unmarshal(response, &res)
	return res
}

func (t *Request) ListDomains(token string) models.Domain {
	body := make(map[string]interface{})
	body["token"] = token
	var res models.Domain

	response := t.post("/domains", body)
	if response == nil {
		return models.Domain{}
	}
	json.Unmarshal(response, &res)
	return res
}

func (t *Request) CreateNewDomain(token string) interface{} {
	var domain string

	fmt.Println("SubDomain:")
	fmt.Scanf("%s", &domain)

	body := make(map[string]interface{})
	body["token"] = token
	body["domain"] = domain
	var res models.Response

	response := t.post("/create_domain", body)
	json.Unmarshal(response, &res)
	return res
}

func (t *Request) RenewDomain(token, domain string) interface{} {

	body := make(map[string]interface{})
	body["token"] = token
	body["domain"] = domain
	var res models.Response

	response := t.post("/renew_domain", body)
	json.Unmarshal(response, &res)
	return res
}

func (t *Request) DeleteDomain(token string, id []string) models.Response {
	body := make(map[string]interface{})
	body["token"] = token
	body["domain_id_list"] = id
	var res models.Response

	response := t.post("/delete_domain", body)
	json.Unmarshal(response, &res)
	return res
}
