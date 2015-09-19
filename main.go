package gha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type token struct {
	Token string `json:"token"`
}

type Request struct {
	Note   string   `json:"note"`
	Scopes []string `json:"scopes"`
}

// Auth gets Psersonal access token of GitHub.
func Auth(user, pass string, r *Request) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	r.Note = fmt.Sprintf("%s for %s@%s", r.Note, os.Getenv("USER"), hostname)
	reqBuf := bytes.NewBuffer([]byte{})
	json.NewEncoder(reqBuf).Encode(r)
	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/authorizations",
		reqBuf,
	)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(user, pass)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		body, _ := ioutil.ReadAll(res.Body)
		return "", fmt.Errorf(string(body))
	}

	t := &token{}
	json.NewDecoder(res.Body).Decode(t)
	return t.Token, nil
}
