package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RobotsAndPencils/go-saml"
)

const (
	host   = "localhost"
	port   = "28080"
	secret = "secret123"
)

var userPassword = "pass123"
var sp saml.ServiceProviderSettings

func init() {
	sp = saml.ServiceProviderSettings{
		PublicCertPath:              "default.crt",
		PrivateKeyPath:              "default.key",
		IDPSSOURL:                   "https://<IDP>/sso/saml",
		IDPSSODescriptorURL:         "http://www.okta.com/<ID>",
		IDPPublicCertPath:           "idp.crt",
		SPSignRequest:               true,
		AssertionConsumerServiceURL: "http://localhost:28080/saml_consume",
	}
	sp.Init()
}

func handleSAMLResponse(w http.ResponseWriter, r *http.Request) {
	encodedXML := r.FormValue("SAMLResponse")

	if encodedXML == "" {
		fmt.Fprintf(w, "SAMLResponse form value missing")
	}

	response, err := saml.ParseEncodedResponse(encodedXML)
	if err != nil {
		fmt.Fprintf(w, "SAMLResponse parse: "+err.Error())
	}

	err = response.Validate(&sp)
	if err != nil {
		fmt.Fprintf(w, "SAMLResponse validation: "+err.Error())
	}

	samlID := response.Assertion.Subject.NameID.Value
	if samlID == "" {
		fmt.Fprintf(w, "SAML attribute identifier uid missing")
	}
	addCookie(w, "auth", samlID)
}

func createAuthRequestURL() string {
	// generate the AuthnRequest and then get a base64 encoded string of the XML
	authnRequest := sp.GetAuthnRequest()
	b64XML, err := authnRequest.EncodedSignedString(sp.PrivateKeyPath)
	if err != nil {
		panic(err)
	}

	// for convenience, get a URL formed with the SAMLRequest parameter
	url, err := saml.GetAuthnRequestURL(sp.IDPSSOURL, b64XML, "")
	if err != nil {
		panic(err)
	}
	return url
}

// HelloHandler handles main page
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "static/form.html")
	default:
		fmt.Fprintf(w, "Sorry, only GET method is supported.")
	}
}

func addCookie(w http.ResponseWriter, token string, user string) {
	log.Println("Adding auth cookie for user " + user)
	expire := time.Now().Add(time.Hour * 24)
	cookie := http.Cookie{
		Name:    "auth",
		Value:   token,
		Expires: expire,
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
}

func grant(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCookie, _ := r.Cookie("auth")
		if authCookie == nil {
			idpURL := createAuthRequestURL()
			http.Redirect(w, r, idpURL, http.StatusSeeOther)
		}
		fn(w, r)
	}
}

func main() {
	http.HandleFunc("/", grant(HelloHandler))
	http.HandleFunc("/saml_consume", handleSAMLResponse)
	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		log.Fatal("error starting http server : ", err)
		return
	}
}
