package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
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

func validUser(userid string) (bool, UserInfo) {
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
	userInfo := UserInfo{}
	body, _ := io.ReadAll(res.Body)
	json.Unmarshal(body, &userInfo)

	fmt.Println(string(body))
	return true, userInfo
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

type UserInfo struct {
	CreatedAt     time.Time `json:"created_at"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	FamilyName    string    `json:"family_name"`
	GivenName     string    `json:"given_name"`
	Identities    []struct {
		Provider   string `json:"provider"`
		UserID     string `json:"user_id"`
		Connection string `json:"connection"`
		IsSocial   bool   `json:"isSocial"`
	} `json:"identities"`
	IdpTenantDomain string    `json:"idp_tenant_domain"`
	Locale          string    `json:"locale"`
	Name            string    `json:"name"`
	Nickname        string    `json:"nickname"`
	Picture         string    `json:"picture"`
	UpdatedAt       time.Time `json:"updated_at"`
	UserID          string    `json:"user_id"`
	LastIP          string    `json:"last_ip"`
	LastLogin       time.Time `json:"last_login"`
	LoginsCount     int       `json:"logins_count"`
}
