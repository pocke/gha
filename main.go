package gha

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type token struct {
	Token string `json:"token"`
}

// Auth gets Psersonal access token of GitHub.
func Auth(user, pass, appName string) (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		"https://api.github.com/authorizations",
		strings.NewReader(fmt.Sprintf(`{"note":"%s@%s"}`, appName, hostname)),
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
