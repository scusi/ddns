// Package ddns implements Dynamic DNS HTTP API used by NoIP, Vollmar and others
package ddns // import "github.com/scusi/ddns"

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Parameters holds all information needed to perfom ddns updates
type Parameters struct {
	URL      string
	Hostname string
	MyIP     string
	Username string
	Password string
}

/*
// Usage:
//
// in your program use it like this:

func main() {
	p := new(ddns.Parameters)
	p.Hostname = "mybox.me"
	// optionally set myip parameter to update an ip that is not the one the request originates from
	// p.MyIP = ""
	p.URL = "https://api.isp.net/ddnsUpdate"
	p.User = "me"
	p.Host = "myPasswd"

	client, err := ddns.NewClient(p)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Update()
	if err != nil {
		log.Fatal(err)
	}

	body, err := ddns.ParseResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", body)
}
*/

// URLFromParameters takes parameters and forms the url used to call the API
func URLFromParameters(p *Parameters) (u *url.URL, err error) {
	u, err = url.Parse(p.URL)
	if err != nil {
		return u, err
	}
	u.User = url.UserPassword(p.Username, p.Password)
	q := u.Query()
	q.Set("hostname", p.Hostname)
	if p.MyIP != "" {
		q.Set("myip", p.MyIP)
	}
	u.RawQuery = q.Encode()
	return u, nil
}

// CallAPI takes the URL and makes a request to the API
func CallAPI(u *url.URL) (body []byte, err error) {
	// setup a http client
	client := &http.Client{}
	// make the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return body, err
	}
	// set user-agent
	req.Header.Set("User-Agent", "ScusiDdns/0.1 - http://github.com/scusi/ddns")

	resp, err := client.Do(req)
	if err != nil {
		return body, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Unexpected Statuscode return: ''%d'\n", resp.StatusCode)
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, nil
}
