package gha

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type token struct {
	Token string `json:"token"`
}

func fileExist(fname string) bool {
	_, err := os.Stat(fname)
	return err == nil
}

// CLI gets Psersonal access token of GitHub.
// username and password are got from STDIN
// And save key to the file.
// If you have the key already at the file, CLI returns this.
func CLI(appName, fname string) (string, error) {
	if fileExist(fname) {
		b, err := ioutil.ReadFile(fname)
		return string(b), err
	}

	user, pass, err := GetUserInfo()
	if err != nil {
		return "", err
	}

	key, err := Auth(user, pass, appName)
	if err != nil {
		return "", err
	}

	if err := ioutil.WriteFile(fname, []byte(key), 0600); err != nil {
		return "", err
	}
	return key, nil
}

func GetUserInfo() (string, string, error) {
	user, err := ReadUsername()
	if err != nil {
		return "", "", err
	}

	pass, err := ReadPassword(user)
	if err != nil {
		return "", "", err
	}

	return user, pass, nil
}

func ReadUsername() (string, error) {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return "", err
	}
	defer tty.Close()

	fmt.Print("username: ")
	sc := bufio.NewScanner(tty)
	sc.Split(bufio.ScanLines)
	sc.Scan()
	return sc.Text(), nil
}

func ReadPassword(user string) (string, error) {
	fmt.Printf("password for %s (never stored): ", user)
	res := make([]byte, 0)

	for {
		v, err := ReadCharAsPassword()
		if err != nil {
			return "", err
		}

		if (v == 127 || v == 8) && len(res) > 0 {
			res = res[:len(res)-1]
			os.Stdout.Write([]byte("\b \b"))
		}

		if v == 13 || v == 10 {
			return string(res), nil
		}

		// C-c or C-d
		if v == 3 || v == 4 {
			return "", fmt.Errorf("Exited by user")
		}

		if v != 0 {
			res = append(res, v)
		}
	}
}

func ReadCharAsPassword() (byte, error) {
	tty, err := os.Open("/dev/tty")
	if err != nil {
		return 0, nil
	}
	defer tty.Close()
	fd := int(tty.Fd())

	if oldState, err := terminal.MakeRaw(fd); err != nil {
		return 0, nil
	} else {
		defer terminal.Restore(fd, oldState)
	}

	var buf [1]byte
	if n, err := syscall.Read(fd, buf[:]); n == 0 || err != nil {
		return 0, err
	}
	return buf[0], nil
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
