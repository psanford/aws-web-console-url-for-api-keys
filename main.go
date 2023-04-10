package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	accessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")

	if accessKeyId == "" {
		log.Fatalf("ENV variable AWS_ACCESS_KEY_ID is not set")
	}
	if secretAccessKey == "" {
		log.Fatalf("ENV variable AWS_SECRET_ACCESS_KEY is not set")
	}
	if sessionToken == "" {
		log.Fatalf("ENV variable AWS_SESSION_TOKEN is not set")
	}

	jsonTxt, err := json.Marshal(map[string]string{
		"sessionId":    accessKeyId,
		"sessionKey":   secretAccessKey,
		"sessionToken": sessionToken,
	})
	if err != nil {
		log.Fatal(err)
	}

	loginURLPrefix := "https://signin.aws.amazon.com/federation"
	req, err := http.NewRequest("GET", loginURLPrefix, nil)
	if err != nil {
		log.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("Action", "getSigninToken")
	q.Add("Session", string(jsonTxt))

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("getSigninToken returned non-200 status: %d", resp.StatusCode)
	}

	var signinTokenResp struct {
		SigninToken string `json:"SigninToken"`
	}

	if err = json.Unmarshal([]byte(body), &signinTokenResp); err != nil {
		log.Fatalf("parse signinTokenResp err: %s", err)
	}

	destination := "https://console.aws.amazon.com/"

	loginURL := fmt.Sprintf(
		"%s?Action=login&Issuer=aws-vault&Destination=%s&SigninToken=%s",
		loginURLPrefix,
		url.QueryEscape(destination),
		url.QueryEscape(signinTokenResp.SigninToken),
	)

	fmt.Println(loginURL)
}
