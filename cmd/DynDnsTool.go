// Implements Dynamic DNS HTTP API used by NoIP, Vollmar and others
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Parameters struct {
	Url      string
	Hostname string
	Username string
	Password string
}

var parameters Parameters

func init() {
	flag.StringVar(&parameters.Url, "url", "", "url to send updates to")
	flag.StringVar(&parameters.Hostname, "host", "", "hostname to update")
	flag.StringVar(&parameters.Username, "user", "", "username for auth")
	flag.StringVar(&parameters.Password, "passwd", "", "password for auth")
}

func main() {
	flag.Parse()
	u, err := urlFromParameters(&parameters)
	if err != nil {
		log.Fatal(err)
	}
	body, err := callApi(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}

func urlFromParameters(p *Parameters) (u *url.URL, err error) {
	u, err = url.Parse(p.Url)
	if err != nil {
		return u, err
	}
	u.User = url.UserPassword(parameters.Username, parameters.Password)
	q := u.Query()
	q.Set("hostname", parameters.Hostname)
	u.RawQuery = q.Encode()
	//log.Printf("Final URL: %s\n", u.String())
	return u, nil
}

func callApi(u *url.URL) (body []byte, err error) {
	// make the request
	resp, err := http.Get(u.String())
	if err != nil {
		return body, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Unexpected Statuscode return: %s\n", resp.StatusCode)
		return body, err
	}
	//log.Printf("Response: %v\n", resp)
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	//log.Printf("Body: %s\n", body)
	return body, nil
}
