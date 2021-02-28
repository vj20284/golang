package main

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	saml2 "github.com/russellhaering/gosaml2"
	"github.com/russellhaering/gosaml2/types"
	dsig "github.com/russellhaering/goxmldsig"
)

func getCertStore() *dsig.MemoryX509CertificateStore {
	res, err := http.Get("https://dev-31010350.okta.com/app/exk6072f30AwzljqY5d6/sso/saml/metadata")
	if err != nil {
		panic(err)
	}

	rawMetadata, err := ioutil.ReadAll(res.Body)
	fmt.Printf("Response %s", rawMetadata)
	if err != nil {
		panic(err)
	}

	metadata := &types.EntityDescriptor{}
	err = xml.Unmarshal(rawMetadata, metadata)
	if err != nil {
		panic(err)
	}

	certStore := dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{},
	}

	for _, kd := range metadata.IDPSSODescriptor.KeyDescriptors {
		for idx, xcert := range kd.KeyInfo.X509Data.X509Certificates {
			if xcert.Data == "" {
				panic(fmt.Errorf("metadata certificate(%d) must not be empty", idx))
			}
			certData, err := base64.StdEncoding.DecodeString(xcert.Data)
			if err != nil {
				panic(err)
			}

			idpCert, err := x509.ParseCertificate(certData)
			if err != nil {
				panic(err)
			}

			certStore.Roots = append(certStore.Roots, idpCert)
		}
	}
	return &certStore
}

func main() {
	// We sign the AuthnRequest with a random key because Okta doesn't seem
	// to verify these.

	sp := &saml2.SAMLServiceProvider{
		IdentityProviderSSOURL:      "https://dev-31010350.okta.com/app/dev-31010350_samlproxytest_1/exk6072f30AwzljqY5d6/sso/saml",
		IdentityProviderIssuer:      "http://www.okta.com/exk6072f30AwzljqY5d6",
		ServiceProviderIssuer:       "http://example.com/saml/acs/example",
		AssertionConsumerServiceURL: "http://localhost:28080/saml_consume",
		AudienceURI:                 "http://localhost:28080/1",
		IDPCertificateStore:         getCertStore(),
		AllowMissingAttributes:      true,
	}

	http.HandleFunc("/saml_consume", func(rw http.ResponseWriter, req *http.Request) {
		err := req.ParseForm()
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		assertionInfo, err := sp.RetrieveAssertionInfo(req.FormValue("SAMLResponse"))
		if err != nil {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.InvalidTime {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		if assertionInfo.WarningInfo.NotInAudience {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		fmt.Fprintf(rw, "NameID: %s\n", assertionInfo.NameID)

		fmt.Fprintf(rw, "Assertions:\n")

		for key, val := range assertionInfo.Values {
			fmt.Fprintf(rw, "  %s: %+v\n", key, val)
		}

		fmt.Fprintf(rw, "\n")

		fmt.Fprintf(rw, "Warnings:\n")
		fmt.Fprintf(rw, "%+v\n", assertionInfo.WarningInfo)
	})

	err := http.ListenAndServe(":28080", nil)
	if err != nil {
		panic(err)
	}
}
