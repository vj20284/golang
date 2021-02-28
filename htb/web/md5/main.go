package main

import (
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

/**

<html>
<head>
<title>emdee five for life</title>
</head>
<body style="background-color:powderblue;">
<h1 align='center'>MD5 encrypt this string</h1><h3 align='center'>FG2RlW8yoYWRZvxlL14X</h3><center><form action="" method="post">
<input type="text" name="hash" placeholder="MD5" align='center'></input>
</br>
<input type="submit" value="Submit"></input>
</form></center>
</body>
</html>
**/

func parseHTML(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error opening page ", err)
		return ""
	}
	defer resp.Body.Close()
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error opening page ", err)
		return ""
	}

	val := doc.Find("h3").Text()
	return val
}

func calculateMD5(input string) string {
	fmt.Println("Input is ", input)
	hash := md5.Sum([]byte(input))
	md5 := hex.EncodeToString(hash[:])
	return md5
}

func postResponse(apiurl string, val string) {
	fmt.Println("Posting ", val)
	form := url.Values{}
	form.Add("hash", val)
	resp, err := http.PostForm(apiurl, form)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Body)
	var res map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&res)
	fmt.Println(res)
}

func main() {
	url := "http://167.99.85.16:31834/"
	val := parseHTML(url)
	postResponse(url, calculateMD5(val))
}
