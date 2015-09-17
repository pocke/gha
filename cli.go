package gha

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

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

	user, pass, err := getUserInfo()
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

func getUserInfo() (string, string, error) {
	user, err := readusername()
	if err != nil {
		return "", "", err
	}

	pass, err := readPassword(user)
	if err != nil {
		return "", "", err
	}

	return user, pass, nil
}

func readusername() (string, error) {
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

func readPassword(user string) (string, error) {
	fmt.Printf("password for %s (never stored): ", user)
	res := make([]byte, 0)

	for {
		v, err := readCharAsPassword()
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

func readCharAsPassword() (byte, error) {
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
