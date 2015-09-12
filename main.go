package gha

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

type token struct {
	Token string `json:"token"`
}

func fileExist(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}

func CLI(appName, fname string) (string, error) {
	if fileExist(fname) {
		b, err := ioutil.ReadFile(fname)
		return string(b), err
	}

	fmt.Print("username: ")

	sc := bufio.NewScanner(os.Stdin)
	sc.Split(bufio.ScanLines)
	sc.Scan()
	user := sc.Text()

	fmt.Printf("password for %s (never stored): ", user)
	pass := string(gopass.GetPasswd())

	key, err := Auth(user, pass, appName)
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(fname, []byte(key), 0600); err != nil {
		return "", err
	}
	return key, nil
}

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
