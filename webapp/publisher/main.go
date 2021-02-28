package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/netSkope/npa_publisher_wizard/networkhelper"
)

// EnforceSSLPinning is true
var EnforceSSLPinning = true

func main() {
	stitcher, tenant, err := discover()
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(stitcher + " " + tenant)
}

func register(registrationToken, tenant, publisherCN string) (publisherRegResult string, err error) {
	// Prepare request
	url := fmt.Sprintf("https://%s/api/publishers/%s", tenant, publisherCN)
	bodyStr := "{ \"registered\": true }"
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer([]byte(bodyStr)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+registrationToken)

	// Send request
	client := getPinnedTLSClient(tenant)
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getPinnedTLSClient(tenant string) *http.Client {
	// The steps to get this Hash
	// 1) Run "echo -n | openssl s_client -connect adminhost:adminport -showcerts"
	// 2) Choose second from the top (this will be the issue of our Admin certificate)
	// 3) Get a public key our of it and get SHA256 for it
	newEdgeIoAdminIssuerPublicKeyHash := []byte{37, 35, 51, 168, 227, 171, 183, 35, 147, 214, 73, 154, 187, 172, 202, 134, 4, 250, 239, 168, 70, 129, 204, 195, 229, 83, 29, 68, 204, 137, 100, 80}
	npaGoSkopeComAdminIssuerPublicKeyHash := []byte{132, 68, 233, 129, 91, 218, 44, 45, 27, 189, 193, 134, 222, 221, 28, 186, 168, 135, 190, 190, 23, 200, 253, 140, 79, 0, 162, 170, 24, 17, 91, 5}

	allAcceptedIssuerCertsHash := [][]byte{newEdgeIoAdminIssuerPublicKeyHash, npaGoSkopeComAdminIssuerPublicKeyHash}
	if EnforceSSLPinning {
		var dialer networkhelper.Dialer
		if strings.HasSuffix(tenant, ".local") {
			// We use a self signed certificate for Stack In The box (and we can't do SSL pinning against it)
			// It's untrusted by the system, so even before SSL pinning kicks in, Go rejects it
			dialer = networkhelper.MakeSSLDialerForSelfSignedCert()
		} else {
			dialer = networkhelper.MakeSSLPinnedDialer(allAcceptedIssuerCertsHash)
		}

		return &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{},
				DialTLS:         dialer,
			},
		}
	} else {
		return &http.Client{}
	}
}

func discover() (discoveryStitcher string, discoveryTenant string, err error) {
	tenant := "ns-2413.nskp-stg.npa.stg.local"
	url := fmt.Sprintf("https://%s/api/discovery", tenant)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	client := getPinnedTLSClient(tenant)
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Parse response
	if resp.StatusCode != 200 {
		return "", "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var discovery map[string]interface{}
	if err := json.Unmarshal(body, &discovery); err != nil {
		return "", "", err
	}

	discoveryStitcher = discovery["stitchers"].([]interface{})[0].(string)
	discoveryTenantIf := discovery["tenant"]

	if discoveryTenantIf == nil {
		discoveryTenant = tenant
	} else {
		discoveryTenant = discoveryTenantIf.(string)
	}

	return discoveryStitcher, discoveryTenant, nil
}
