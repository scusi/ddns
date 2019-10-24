package ddns

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	userAgent  string
	Parameters *Parameters
}

func NewClient(p *Parameters) *Client {
	client := &Client{
		/*
			Parameters: Parameters{
				Hostname: hostname
				MyIP: myip
				Username: username
				Password: password
				URL: "https://dynupdate.no-ip.com/nic/update"
			}
		*/
		Parameters: p,
		userAgent:  "ScusiDdns/0.1 - github.com/scusi/ddns",
	}
	return client
}

func (c *Client) UpdateHostname() (r *http.Response, err error) {
	u, err := URLFromParameters(c.Parameters)
	if err != nil {
		return r, err
	}
	// make the request
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return r, err
	}
	// set user-agent
	req.Header.Set("User-Agent", c.userAgent)
	r, err = c.httpClient.Do(req)
	if err != nil {
		return r, err
	}
	return r, nil
}

func ParseResponse(resp *http.Response) (body []byte, err error) {
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

func (c *Client) Offline() {}

/*
Responses as defind by NoIP.com
good IP_ADDRESS		Success	DNS hostname update successful.
					Followed by a space and the IP address it was updated to.
nochg IP_ADDRESS	Success	IP address is current, no update performed.
					Followed by a space and the IP address that it is currently set to.
nohost		Error	Hostname supplied does not exist under specified account, client exit and require user to enter new login credentials before performing an additional request.
badauth		Error	Invalid username password combination
badagent	Error	Client disabled. Client should exit and not perform any more updates without user intervention.
!donator	Error	An update request was sent including a feature that is not available to that particular user such as offline options.
abuse		Error	Username is blocked due to abuse.
					Either for not following our update specifications or disabled due to violation of the No-IP terms of service.
					Our terms of service can be viewed here.
					Client should stop sending updates.
911	Error			A fatal error on our side such as a database outage. Retry the update no sooner than 30 minutes.
*/
