// Implements Dynamic DNS HTTP API used by NoIP, Vollmar and others
package ddns // import "github.com/scusi/DynDNSClient/lib"

import (
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

func init() {
}

/*
// Usage:
//
// in your program use it like this:

func main() {
	p := new(ddns.Parameters)
	p.Hostname = "mybox.me"
	p. URL = "https://api.isp.net/ddnsUpdate"
	p.User = "me"
	p.Host = "myPasswd"
	u, err := ddns.UrlFromParameters(&parameters)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ddns.CallApi(u)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)
}
*/

func UrlFromParameters(p *Parameters) (u *url.URL, err error) {
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

func CallApi(u *url.URL) (body []byte, err error) {
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
