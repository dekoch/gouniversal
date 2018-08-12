package getPublicIP

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

var servers = [...]string{"https://bot.whatismyipaddress.com",
	"https://ident.me",
	"https://api.ipify.org?format=text",
	"https://myexternalip.com/raw"}

func Get() (string, error) {

	for _, server := range servers {

		fmt.Print("get ip from " + server + "...")

		resp, err := http.Get(server)
		if err != nil {
			fmt.Println(err)
			continue
		} else {
			defer resp.Body.Close()
		}

		ipresp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ipstr := string(ipresp)
		ipstr = strings.TrimSpace(ipstr)

		ip := net.ParseIP(ipstr)

		fmt.Println(ip.String())

		if ip.To4() != nil {

			return ip.String(), nil

		} else if ip.To16() != nil {

			return ip.String(), nil
		}
	}

	return "", errors.New("no ip found")
}
