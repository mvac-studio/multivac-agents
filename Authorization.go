package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type Credentials struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type CredentialsRequest struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

func validUser(userid string) bool {
	domain := os.Getenv("AUTH_DOMAIN")
	clientId := os.Getenv("AUTH_CLIENT_ID")
	secret := os.Getenv("AUTH_CLIENT_SECRET")

	fmt.Printf("domain: %s\nclientId: %s\nsecret: %s\n", domain, clientId, secret)
	url := fmt.Sprintf("https://%s/api/v2/users/%s", domain, userid)

	credentials := getCredentials(domain, clientId, secret)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", credentials.AccessToken))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	return true
}

func getCredentials(domain string, clientId string, secret string) *Credentials {
	url := fmt.Sprintf("https://%s/oauth/token", domain)

	request := CredentialsRequest{
		ClientId:     clientId,
		ClientSecret: secret,
		Audience:     fmt.Sprintf("https://%s/api/v2/", domain),
		GrantType:    "client_credentials",
	}

	data, err := json.Marshal(request)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(data)))

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	body, _ := io.ReadAll(res.Body)
	credentials := Credentials{}
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		fmt.Println(err)
	}
	return &credentials
}
